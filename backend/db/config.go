package db

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Name            string        `yaml:"name"`
	ExtraParams     string        `yaml:"extra_params"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

// 从YAML文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config struct {
		Database Config `yaml:"database"`
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config.Database, nil
}

func (c Config) DSN() string {
	host := c.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := c.Port
	if port == 0 {
		port = 3306
	}
	name := c.Name
	if name == "" {
		name = "smart_forest"
	}
	params := "charset=utf8mb4&parseTime=true&loc=Local"
	if c.ExtraParams != "" {
		params = c.ExtraParams
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		c.User,
		c.Password,
		host,
		port,
		name,
		params,
	)
}
