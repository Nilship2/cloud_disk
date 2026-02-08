// internal/handler/router.go
package handler

import (
	"cloud-disk/internal/config"
	"cloud-disk/internal/handler/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config) *gin.Engine {
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

	// API v1 路由组
	//v1 := router.Group("/api/v1")
	//{
	// 用户相关路由（第二天实现）
	// v1.POST("/register", userHandler.Register)
	// v1.POST("/login", userHandler.Login)

	// 需要认证的路由组
	// authorized := v1.Group("/")
	// authorized.Use(middleware.AuthRequired())
	// {
	//     authorized.GET("/profile", userHandler.GetProfile)
	//     authorized.PUT("/profile", userHandler.UpdateProfile)
	// }
	//}

	return router
}
