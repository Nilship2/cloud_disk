// internal/config/load.go
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

var GlobalConfig *Config

func Load(configPath ...string) (*Config, error) {
	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 设置配置文件搜索路径
	if len(configPath) > 0 && configPath[0] != "" {
		v.SetConfigFile(configPath[0])
	} else {
		// 从当前目录和config目录查找
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/cloud-disk/")
	}

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, using defaults")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// 读取环境变量（前缀为CLOUD_DISK_）
	v.SetEnvPrefix("CLOUD_DISK")
	v.AutomaticEnv()

	// 绑定环境变量
	bindEnvVars(v)

	// 解析配置到结构体
	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	GlobalConfig = config
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.Username = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		config.Database.Database = name
	}

	return config, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.env", "development")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)

	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.username", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.database", "cloud_disk")
	v.SetDefault("database.charset", "utf8mb4")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.conn_max_lifetime", 3600)

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 20)

	v.SetDefault("storage.type", "local")
	v.SetDefault("storage.local.base_path", "./storage/uploads")
	v.SetDefault("storage.local.temp_path", "./storage/uploads/temp")
	v.SetDefault("storage.local.max_size_mb", 100)

	v.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	v.SetDefault("jwt.expires_hours", 24)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")
}

func bindEnvVars(v *viper.Viper) {
	// 服务器相关
	v.BindEnv("server.env", "CLOUD_DISK_SERVER_ENV")
	v.BindEnv("server.port", "CLOUD_DISK_SERVER_PORT")

	// 数据库相关
	v.BindEnv("database.host", "CLOUD_DISK_DB_HOST")
	v.BindEnv("database.port", "CLOUD_DISK_DB_PORT")
	v.BindEnv("database.username", "CLOUD_DISK_DB_USERNAME")
	v.BindEnv("database.password", "CLOUD_DISK_DB_PASSWORD")
	v.BindEnv("database.database", "CLOUD_DISK_DB_NAME")

	// JWT
	v.BindEnv("jwt.secret", "CLOUD_DISK_JWT_SECRET")

	// 存储
	v.BindEnv("storage.type", "CLOUD_DISK_STORAGE_TYPE")
}

func validateConfig(config *Config) error {
	// 验证服务器配置
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// 只在生产环境验证JWT密钥
	if config.Server.Env == "production" {
		if config.JWT.Secret == "" ||
			config.JWT.Secret == "your-secret-key-change-in-production" {
			return fmt.Errorf("JWT secret must be set in production")
		}
	}

	// 确保存储目录存在
	if config.Storage.Type == "local" {
		basePath := config.Storage.Local.BasePath
		if err := ensureDirectoryExists(basePath); err != nil {
			return fmt.Errorf("failed to create storage directory: %w", err)
		}

		tempPath := config.Storage.Local.TempPath
		if err := ensureDirectoryExists(tempPath); err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
	}

	return nil
}

func ensureDirectoryExists(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	return os.MkdirAll(absPath, 0755)
}
