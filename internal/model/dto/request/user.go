// internal/model/dto/request/user.go
package request

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 注册请求
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50,username"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=32"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 更新用户信息请求
type UpdateProfileRequest struct {
	Avatar string `json:"avatar" binding:"omitempty,url"`
	Bio    string `json:"bio" binding:"omitempty,max=500"`
}

// 修改密码请求
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=32"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// 邮箱验证请求
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// ValidateUsername 自定义验证器：用户名只能包含字母、数字和下划线
func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}
