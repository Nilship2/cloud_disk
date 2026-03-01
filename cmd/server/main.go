// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"cloud-disk/internal/auth"
	"cloud-disk/internal/config"
	"cloud-disk/internal/database"
	"cloud-disk/internal/handler"
	"cloud-disk/pkg/logger"
)

// @title 云盘系统 API
// @version 1.0
// @description 仿百度网盘的后端系统，支持文件存储、分享、收藏等功能
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// 1. 确定环境
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// 2. 根据环境选择配置文件
	configFile := "config.dev.yaml"
	if env == "production" {
		configFile = "config.prod.yaml"
	} else if env == "test" {
		configFile = "config.test.yaml"
	}

	log.Printf("Running in %s mode, using config file: %s", env, configFile)

	// 3. 加载配置
	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 4. 从环境变量覆盖数据库配置（Docker部署时使用）
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
		log.Printf("Using DB_HOST from environment: %s", host)
	}
	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Database.Port = port
			log.Printf("Using DB_PORT from environment: %d", port)
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.Username = user
		log.Printf("Using DB_USER from environment: %s", user)
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
		log.Printf("Using DB_PASSWORD from environment")
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.Database.Database = name
		log.Printf("Using DB_NAME from environment: %s", name)
	}

	// 5. 从环境变量覆盖 JWT 配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}

	// 6. 从环境变量覆盖 MinIO 配置
	if accessKey := os.Getenv("MINIO_ACCESS_KEY"); accessKey != "" {
		cfg.Storage.MinIO.AccessKey = accessKey
	}
	if secretKey := os.Getenv("MINIO_SECRET_KEY"); secretKey != "" {
		cfg.Storage.MinIO.SecretKey = secretKey
	}

	// 2. 初始化日志
	if err := logger.Init(&cfg.Log); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Configuration loaded successfully")
	logger.Infof("Server starting in %s mode on %s:%d",
		cfg.Server.Env, cfg.Server.Host, cfg.Server.Port)

	// 3. 初始化数据库连接
	db, err := database.Init(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 4. 自动迁移数据库表结构（开发环境）
	if cfg.Server.Env == "development" {
		if err := database.AutoMigrate(db); err != nil {
			logger.Warn("Auto migration failed", logger.Any("error", err))
		}
	}

	// 5. 初始化JWT管理器
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiresHours)

	// 6. 初始化路由
	router := handler.NewRouter(cfg, db, jwtManager)

	// 7. 创建HTTP服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.Server.WriteTimeout * time.Second,
	}

	// 8. 启动服务器
	go func() {
		logger.Infof("Server listening on http://%s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 9. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 10. 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited properly")
}
