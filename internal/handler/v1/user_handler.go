// internal/handler/v1/user_handler.go
package v1

import (
	"cloud-disk/internal/constant"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService interfaces.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService interfaces.UserService) *UserHandler {
	validate := validator.New()

	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("username", request.ValidateUsername)
		validate = v
	} else {
		// 后备方案
		validate.RegisterValidation("username", request.ValidateUsername)
	}

	return &UserHandler{
		userService: userService,
		validator:   validate,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户，用户名唯一，邮箱唯一
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=entity.UserResponse} "注册成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 409 {object} response.ErrorResponse "用户已存在"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	// 参数验证
	if err := h.validator.Struct(req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrUserExists, err.Error())
		return
	}

	response.Created(c, user.ToResponse())
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录，返回JWT令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=response.LoginResponse} "登录成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 401 {object} response.ErrorResponse "密码错误"
// @Failure 404 {object} response.ErrorResponse "用户不存在"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	// 参数验证
	if err := h.validator.Struct(req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	loginResp, err := h.userService.Login(&req)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrPasswordWrong, err.Error())
		return
	}

	response.Success(c, loginResp)
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=entity.UserResponse} "成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrUserNotExists, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新个人信息
// @Summary 更新个人信息
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body request.UpdateProfileRequest true "个人信息"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Router /api/v1/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.userService.UpdateProfile(userID, &req); err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, nil)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body request.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Router /api/v1/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.userService.ChangePassword(userID, &req); err != nil {
		response.ErrorWithMessage(c, constant.ErrPasswordWrong, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetStorageInfo 获取存储空间信息
// @Summary 获取存储空间信息
// @Tags 用户管理
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=response.StorageInfo} "成功"
// @Router /api/v1/storage [get]
func (h *UserHandler) GetStorageInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	storageInfo, err := h.userService.GetStorageInfo(userID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, storageInfo)
}
