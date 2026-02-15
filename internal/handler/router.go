// internal/handler/router.go
package handler

import (
	"net/http"

	"cloud-disk/internal/auth"
	"cloud-disk/internal/config"
	"cloud-disk/internal/dao"
	"cloud-disk/internal/handler/middleware"
	v1 "cloud-disk/internal/handler/v1"
	"cloud-disk/internal/service/impl"
	storageImpl "cloud-disk/pkg/storage" // 使用别名

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
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"env":    cfg.Server.Env,
		})
	})

	// 静态文件服务（本地存储）
	if cfg.Storage.Type == "local" {
		router.Static("/files", cfg.Storage.Local.BasePath)
	}

	// 初始化DAO
	userDAO := dao.NewUserDAO(db)
	fileDAO := dao.NewFileDAO(db)
	shareDAO := dao.NewShareDAO(db)
	favoriteDAO := dao.NewFavoriteDAO(db)

	// 初始化存储
	var storage storageImpl.Storage
	if cfg.Storage.Type == "local" {
		storage = storageImpl.NewLocalStorage(
			cfg.Storage.Local.BasePath,
			cfg.Storage.Local.TempPath,
			int64(cfg.Storage.Local.MaxSizeMB),
		)
	}
	// TODO: 添加 MinIO 和 S3 支持

	// 初始化服务
	userService := impl.NewUserService(db, jwtManager)
	fileService := impl.NewFileService(db, fileDAO, userDAO, storage)
	shareService := impl.NewShareService(db, shareDAO, fileDAO, userDAO, storage) // 现在 storage 已定义
	favoriteService := impl.NewFavoriteService(db, favoriteDAO, fileDAO)

	// 初始化处理器
	userHandler := v1.NewUserHandler(userService)
	fileHandler := v1.NewFileHandler(fileService)
	shareHandler := v1.NewShareHandler(shareService)
	favoriteHandler := v1.NewFavoriteHandler(favoriteService)

	// API v1 路由组
	v1Group := router.Group("/api/v1")
	{
		// 公开路由
		v1Group.POST("/register", userHandler.Register)
		v1Group.POST("/login", userHandler.Login)

		// 公开分享路由（不需要认证）
		v1Group.GET("/s/:token", shareHandler.Access)
		v1Group.GET("/s/:token/download", shareHandler.Download)

		// 需要认证的路由组
		authorized := v1Group.Group("/")
		authorized.Use(middleware.AuthRequired(jwtManager))
		{
			// 用户相关
			authorized.GET("/profile", userHandler.GetProfile)
			authorized.PUT("/profile", userHandler.UpdateProfile)
			authorized.POST("/change-password", userHandler.ChangePassword)
			authorized.GET("/storage", userHandler.GetStorageInfo)

			// 文件相关
			authorized.POST("/files/upload", fileHandler.Upload)
			authorized.POST("/files/instant", fileHandler.InstantUpload)
			authorized.GET("/files", fileHandler.GetList)
			authorized.GET("/files/check", fileHandler.CheckExists)
			authorized.GET("/files/:id", fileHandler.GetDetail)
			authorized.GET("/files/:id/download", fileHandler.Download)
			authorized.DELETE("/files/:id", fileHandler.Delete)
			authorized.DELETE("/files/batch", fileHandler.BatchDelete)
			authorized.PUT("/files/:id/rename", fileHandler.Rename)
			authorized.PUT("/files/:id/move", fileHandler.Move)

			// 文件夹相关
			authorized.POST("/folders", fileHandler.CreateFolder)

			// 分享相关
			authorized.POST("/shares", shareHandler.Create)
			authorized.GET("/shares", shareHandler.GetList)
			authorized.GET("/shares/:id", shareHandler.GetDetail)
			authorized.PUT("/shares/:id", shareHandler.Update)
			authorized.DELETE("/shares/:id", shareHandler.Cancel)

			// 收藏相关
			authorized.POST("/favorites", favoriteHandler.Add)
			authorized.GET("/favorites", favoriteHandler.GetList)
			authorized.GET("/favorites/check", favoriteHandler.Check)
			authorized.DELETE("/favorites/:file_id", favoriteHandler.Remove)
		}
	}

	return router
}
