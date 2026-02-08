// pkg/response/error_response.go
package response

import "github.com/gin-gonic/gin"

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

func ValidationError(c *gin.Context, errors []ErrorDetail) {
	c.JSON(400, ErrorResponse{
		Code:    1001,
		Message: "参数验证失败",
		Errors:  errors,
	})
}
