// internal/handler/middleware/auth.go
package middleware

import (
	"strings"

	"cloud-disk/internal/auth"
	"cloud-disk/internal/constant"
	"cloud-disk/pkg/logger"
	"cloud-disk/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthRequired JWT认证中间件
func AuthRequired(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取token
		token := extractToken(c)
		if token == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 2. 验证token
		claims, err := jwtManager.Validate(token)
		if err != nil {
			logger.Error("Token validation failed", zap.Any("error", err))

			switch err {
			case auth.ErrTokenExpired:
				response.ErrorWithMessage(c, constant.ErrUnauthorized, "token已过期")
			case auth.ErrTokenMalformed:
				response.ErrorWithMessage(c, constant.ErrUnauthorized, "token格式错误")
			default:
				response.Unauthorized(c)
			}
			c.Abort()
			return
		}

		// 3. 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// extractToken 从请求中提取token
func extractToken(c *gin.Context) string {
	// 1. 从Authorization header获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// 2. 从查询参数获取
	token := c.Query("token")
	if token != "" {
		return token
	}

	// 3. 从Cookie获取
	token, _ = c.Cookie("token")
	return token
}

// GetCurrentUserID 获取当前用户ID
func GetCurrentUserID(c *gin.Context) uint {
	return c.GetUint("user_id")
}

// GetCurrentUsername 获取当前用户名
func GetCurrentUsername(c *gin.Context) string {
	return c.GetString("username")
}

// GetCurrentEmail 获取当前用户邮箱
func GetCurrentEmail(c *gin.Context) string {
	return c.GetString("email")
}
