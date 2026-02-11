// internal/handler/router.go
package handler

import (
	"cloud-disk/internal/auth"
	"cloud-disk/internal/config"
	"cloud-disk/internal/handler/middleware"
	v1 "cloud-disk/internal/handler/v1"
	"cloud-disk/internal/service/impl"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, db *gorm.DB, jwtManager *auth.JWTManager) *gin.Engine {
	// 设置Gin模式
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// 全局中间件
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.Recovery())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    cfg.Server.Env,
		})
	})

	// 初始化服务
	userService := impl.NewUserService(db, jwtManager)

	// 初始化处理器
	userHandler := v1.NewUserHandler(userService)

	// API v1 路由组
	v1Group := router.Group("/api/v1")
	{
		// 公开路由
		v1Group.POST("/register", userHandler.Register)
		v1Group.POST("/login", userHandler.Login)

		// 需要认证的路由组
		authorized := v1Group.Group("/")
		authorized.Use(middleware.AuthRequired(jwtManager))
		{
			// 用户相关
			authorized.GET("/profile", userHandler.GetProfile)
			authorized.PUT("/profile", userHandler.UpdateProfile)
			authorized.POST("/change-password", userHandler.ChangePassword)
			authorized.GET("/storage", userHandler.GetStorageInfo)
		}
	}

	return router
}
