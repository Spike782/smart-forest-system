package util

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"smart-forest-system/backend/db"
)

// BackupConfig 备份配置
type BackupConfig struct {
	BackupDir      string        `yaml:"backup_dir"`      // 备份文件存储目录
	BackupInterval time.Duration `yaml:"backup_interval"` // 备份间隔
	RetentionDays  int           `yaml:"retention_days"`  // 备份保留天数
	DBConfig       db.Config     `yaml:"db_config"`       // 数据库配置
}

// BackupResult 备份结果
type BackupResult struct {
	BackupFile string    `json:"backup_file"` // 备份文件路径
	Size       int64     `json:"size"`       // 备份文件大小
	Duration   float64   `json:"duration"`   // 备份持续时间（秒）
	Timestamp  time.Time `json:"timestamp"`  // 备份时间
}

// BackupDatabase 备份数据库
func BackupDatabase(db *sql.DB, config BackupConfig) (*BackupResult, error) {
	// 确保备份目录存在
	if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
		return nil, fmt.Errorf("创建备份目录失败: %v", err)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(config.BackupDir, fmt.Sprintf("smart_forest_backup_%s.sql", timestamp))

	// 开始备份
	startTime := time.Now()

	// 使用 mysqldump 命令备份数据库
	// 注意：在实际环境中，应该使用配置文件中的数据库连接信息
	dsn := config.DBConfig.DSN()
	cmd := fmt.Sprintf("mysqldump --single-transaction --routines --triggers --events %s > %s", dsn, backupFile)

	log.Printf("执行备份命令: %s", cmd)

	// 执行备份命令
	if err := runCommand(cmd); err != nil {
		return nil, fmt.Errorf("执行备份命令失败: %v", err)
	}

	// 计算备份时间和文件大小
	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()

	// 获取备份文件大小
	fileInfo, err := os.Stat(backupFile)
	if err != nil {
		return nil, fmt.Errorf("获取备份文件信息失败: %v", err)
	}

	// 清理过期备份
	if err := cleanupOldBackups(config); err != nil {
		log.Printf("清理过期备份失败: %v", err)
	}

	return &BackupResult{
		BackupFile: backupFile,
		Size:       fileInfo.Size(),
		Duration:   duration,
		Timestamp:  endTime,
	}, nil
}

// RestoreDatabase 恢复数据库
func RestoreDatabase(db *sql.DB, backupFile string, config BackupConfig) error {
	// 检查备份文件是否存在
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", backupFile)
	}

	// 使用 mysql 命令恢复数据库
	dsn := config.DBConfig.DSN()
	cmd := fmt.Sprintf("mysql %s < %s", dsn, backupFile)

	log.Printf("执行恢复命令: %s", cmd)

	// 执行恢复命令
	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("执行恢复命令失败: %v", err)
	}

	return nil
}

// ScheduleBackup 定时备份数据库
func ScheduleBackup(config BackupConfig) {
	// 初始化数据库连接
	database := db.InitDB()
	defer database.Close()

	// 创建备份目录
	if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
		log.Fatalf("创建备份目录失败: %v", err)
	}

	// 立即执行一次备份
	if _, err := BackupDatabase(database, config); err != nil {
		log.Printf("立即备份失败: %v", err)
	}

	// 定时执行备份
	ticker := time.NewTicker(config.BackupInterval)
	defer ticker.Stop()

	log.Printf("开始定时备份，间隔: %v", config.BackupInterval)

	for {
		select {
		case <-ticker.C:
			if _, err := BackupDatabase(database, config); err != nil {
				log.Printf("定时备份失败: %v", err)
			} else {
				log.Printf("定时备份成功")
			}
		}
	}
}

// cleanupOldBackups 清理过期备份
func cleanupOldBackups(config BackupConfig) error {
	// 获取备份目录中的所有文件
	files, err := os.ReadDir(config.BackupDir)
	if err != nil {
		return fmt.Errorf("读取备份目录失败: %v", err)
	}

	// 计算过期时间
	expiryTime := time.Now().AddDate(0, 0, -config.RetentionDays)

	// 遍历所有文件，删除过期备份
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 获取文件信息
		fileInfo, err := file.Info()
		if err != nil {
			log.Printf("获取文件信息失败: %v", err)
			continue
		}

		// 检查文件是否过期
		if fileInfo.ModTime().Before(expiryTime) {
			filePath := filepath.Join(config.BackupDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				log.Printf("删除过期备份失败: %v", err)
			} else {
				log.Printf("删除过期备份: %s", filePath)
			}
		}
	}

	return nil
}

// runCommand 执行系统命令
func runCommand(cmd string) error {
	// 使用 PowerShell 执行命令（因为是 Windows 环境）
	cmdArgs := []string{"-Command", cmd}
	proc, err := os.StartProcess("powershell.exe", cmdArgs, &os.ProcAttr{
		Dir:   ".",
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})

	if err != nil {
		return fmt.Errorf("启动命令失败: %v", err)
	}

	// 等待命令执行完成
	state, err := proc.Wait()
	if err != nil {
		return fmt.Errorf("等待命令执行完成失败: %v", err)
	}

	if !state.Success() {
		return fmt.Errorf("命令执行失败，退出码: %d", state.ExitCode())
	}

	return nil
}
