// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud-disk/internal/auth"
	"cloud-disk/internal/config"
	"cloud-disk/internal/database"
	"cloud-disk/internal/handler"
	"cloud-disk/pkg/logger"
)

func main() {
	// 1. 加载配置 - 使用 SQLite 配置
	cfg, err := config.Load("config.mysql.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
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
