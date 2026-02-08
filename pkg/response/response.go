// pkg/response/response.go
package response

import (
	"cloud-disk/internal/constant"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(constant.StatusSuccess, Response{
		Code:    constant.SuccessCode,
		Message: constant.GetErrorMessage(constant.SuccessCode),
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(constant.StatusSuccess, Response{
		Code:    constant.SuccessCode,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(constant.StatusCreated, Response{
		Code:    constant.SuccessCode,
		Message: "创建成功",
		Data:    data,
	})
}

func Error(c *gin.Context, code int) {
	c.JSON(getHttpStatus(code), Response{
		Code:    code,
		Message: constant.GetErrorMessage(code),
		Error:   constant.GetErrorMessage(code),
	})
}

func ErrorWithMessage(c *gin.Context, code int, message string) {
	c.JSON(getHttpStatus(code), Response{
		Code:    code,
		Message: message,
		Error:   constant.GetErrorMessage(code),
	})
}

func BadRequest(c *gin.Context, code int) {
	c.JSON(constant.StatusBadRequest, Response{
		Code:    code,
		Message: constant.GetErrorMessage(code),
		Error:   constant.GetErrorMessage(code),
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(constant.StatusUnauthorized, Response{
		Code:    constant.ErrUnauthorized,
		Message: constant.GetErrorMessage(constant.ErrUnauthorized),
		Error:   constant.GetErrorMessage(constant.ErrUnauthorized),
	})
}

func Forbidden(c *gin.Context) {
	c.JSON(constant.StatusForbidden, Response{
		Code:    constant.ErrPermissionDenied,
		Message: constant.GetErrorMessage(constant.ErrPermissionDenied),
		Error:   constant.GetErrorMessage(constant.ErrPermissionDenied),
	})
}

func NotFound(c *gin.Context) {
	c.JSON(constant.StatusNotFound, Response{
		Code:    constant.ErrNotFound,
		Message: constant.GetErrorMessage(constant.ErrNotFound),
		Error:   constant.GetErrorMessage(constant.ErrNotFound),
	})
}

func InternalError(c *gin.Context) {
	c.JSON(constant.StatusInternalServerError, Response{
		Code:    constant.ErrInternal,
		Message: constant.GetErrorMessage(constant.ErrInternal),
		Error:   constant.GetErrorMessage(constant.ErrInternal),
	})
}

func getHttpStatus(code int) int {
	// 根据错误码映射HTTP状态码
	switch {
	case code >= 2000 && code < 3000:
		return constant.StatusBadRequest
	case code >= 3000 && code < 4000:
		return constant.StatusBadRequest
	case code >= 4000 && code < 5000:
		return constant.StatusBadRequest
	default:
		return constant.StatusInternalServerError
	}
}
