package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 初始化数据库连接
func InitDB() *sql.DB {
	return InitDBWithConfig("db/config.yaml")
}

func InitDBWithConfig(configPath string) *sql.DB {
	fmt.Printf("正在加载配置文件: %s\n", configPath)

	// 加载配置
	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载数据库配置失败: %v", err)
	}
	fmt.Printf("配置加载成功，DSN: %s\n", cfg.DSN())

	// 打开数据库连接
	db, err := OpenDatabase(*cfg)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	fmt.Println("数据库连接成功!")
	return db
}

func OpenDatabase(cfg Config) (*sql.DB, error) {
	dsn := cfg.DSN()
	pool, err := sql.Open("mysql", dsn)
	if err != nil {

		return nil, err
	}

	// 设置数据库连接池的最大打开连接数
	if cfg.MaxOpenConns > 0 {
		pool.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	// 设置数据库连接池的最大空闲连接数
	if cfg.MaxIdleConns > 0 {
		pool.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	// 设置连接的最大生命周期
	if cfg.ConnMaxLifetime > 0 {
		pool.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	// 设置连接的最大空闲时间
	if cfg.ConnMaxIdleTime > 0 {
		pool.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	// 使用5秒超时上下文测试数据库连接是否有效
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.PingContext(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
