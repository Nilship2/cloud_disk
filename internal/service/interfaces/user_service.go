// internal/service/interfaces/user_service.go
package interfaces

import (
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
	"cloud-disk/internal/model/entity"
)

type UserService interface {
	// 注册
	Register(req *request.RegisterRequest) (*entity.User, error)

	// 登录
	Login(req *request.LoginRequest) (*response.LoginResponse, error)

	// 获取用户信息
	GetProfile(userID uint) (*entity.UserResponse, error)

	// 更新个人信息
	UpdateProfile(userID uint, req *request.UpdateProfileRequest) error

	// 修改密码
	ChangePassword(userID uint, req *request.ChangePasswordRequest) error

	// 获取存储空间信息
	GetStorageInfo(userID uint) (*response.StorageInfo, error)

	// 检查存储空间
	CheckStorage(userID uint, fileSize int64) (bool, error)

	// 增加已用空间
	AddUsedSpace(userID uint, size int64) error

	// 减少已用空间
	ReduceUsedSpace(userID uint, size int64) error
}
