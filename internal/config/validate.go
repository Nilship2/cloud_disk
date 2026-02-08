// internal/config/validate.go
package config

import (
	"errors"
	"strings"
)

func Validate() error {
	if GlobalConfig == nil {
		return errors.New("config not loaded")
	}

	// 验证环境变量
	validEnvs := map[string]bool{
		"development": true,
		"production":  true,
		"test":        true,
	}

	if !validEnvs[GlobalConfig.Server.Env] {
		return errors.New("invalid environment, must be one of: development, production, test")
	}

	// 验证存储类型
	validStorageTypes := map[string]bool{
		"local": true,
		"minio": true,
		"s3":    true,
	}

	if !validStorageTypes[GlobalConfig.Storage.Type] {
		return errors.New("invalid storage type, must be one of: local, minio, s3")
	}

	// 验证日志级别
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}

	logLevel := strings.ToLower(GlobalConfig.Log.Level)
	if !validLogLevels[logLevel] {
		return errors.New("invalid log level")
	}

	return nil
}
