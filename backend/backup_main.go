package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"smart-forest-system/backend/util"

	"gopkg.in/yaml.v2"
)

func main() {
	// 解析命令行参数
	var (  
		configFile string
		backup     bool
		restore    bool
		backupFile string
	)

	flag.StringVar(&configFile, "config", "./config/backup.yaml", "备份配置文件路径")
	flag.BoolVar(&backup, "backup", false, "执行数据库备份")
	flag.BoolVar(&restore, "restore", false, "执行数据库恢复")
	flag.StringVar(&backupFile, "file", "", "恢复使用的备份文件路径")
	flag.Parse()

	// 加载配置文件
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 根据命令执行相应操作
	if backup {
		// 执行备份
		result, err := util.BackupDatabase(nil, config)
		if err != nil {
			log.Fatalf("备份数据库失败: %v", err)
		}
		fmt.Printf("备份成功:\n")
		fmt.Printf("  备份文件: %s\n", result.BackupFile)
		fmt.Printf("  文件大小: %d 字节\n", result.Size)
		fmt.Printf("  备份时间: %v\n", result.Timestamp)
		fmt.Printf("  持续时间: %.2f 秒\n", result.Duration)
	} else if restore {
		// 执行恢复
		if backupFile == "" {
			log.Fatalf("恢复操作必须指定备份文件路径")
		}
		if err := util.RestoreDatabase(nil, backupFile, config); err != nil {
			log.Fatalf("恢复数据库失败: %v", err)
		}
		fmt.Printf("恢复成功: %s\n", backupFile)
	} else {
		// 启动定时备份服务
		fmt.Println("启动定时备份服务...")
		fmt.Printf("备份间隔: %v\n", config.BackupInterval)
		fmt.Printf("备份目录: %s\n", config.BackupDir)
		fmt.Printf("保留天数: %d\n", config.RetentionDays)
		fmt.Println("按 Ctrl+C 停止服务")
		util.ScheduleBackup(config)
	}
}

// loadConfig 加载备份配置
func loadConfig(configFile string) (util.BackupConfig, error) {
	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return util.BackupConfig{}, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析 YAML 配置
	var config util.BackupConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return util.BackupConfig{}, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if config.BackupDir == "" {
		// 如果未指定备份目录，使用当前目录下的 backups 子目录
		config.BackupDir = filepath.Join(".", "backups")
	}

	if config.BackupInterval == 0 {
		// 默认备份间隔为 24 小时
		config.BackupInterval = 24 * time.Hour
	}

	if config.RetentionDays == 0 {
		// 默认保留 7 天备份
		config.RetentionDays = 7
	}

	return config, nil
}
