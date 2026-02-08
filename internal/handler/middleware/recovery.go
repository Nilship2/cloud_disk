// internal/handler/middleware/recovery.go
package middleware

import (
	"cloud-disk/pkg/logger"
	"cloud-disk/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", zap.Any("error", err))

				// 使用统一的错误响应
				response.InternalError(c)
				c.Abort()
			}
		}()

		c.Next()
	}
}
