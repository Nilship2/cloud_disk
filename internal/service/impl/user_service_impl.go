// internal/service/impl/user_service_impl.go
package impl

import (
	"errors"

	"cloud-disk/internal/auth"
	"cloud-disk/internal/constant"
	"cloud-disk/internal/dao"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
	"cloud-disk/internal/model/entity"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/internal/utils/crypto"
	"cloud-disk/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	userDAO    *dao.UserDAO
	jwtManager *auth.JWTManager
}

func NewUserService(db *gorm.DB, jwtManager *auth.JWTManager) interfaces.UserService {
	return &UserServiceImpl{
		userDAO:    dao.NewUserDAO(db),
		jwtManager: jwtManager,
	}
}

// Register 用户注册
func (s *UserServiceImpl) Register(req *request.RegisterRequest) (*entity.User, error) {
	// 1. 检查用户是否已存在
	exists, err := s.userDAO.Exists(req.Username, req.Email)
	if err != nil {
		logger.Error("Failed to check user existence", logger.Any("error", err))
		return nil, errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}
	if exists {
	} //防止变量未引用的报错，以后可能要记得删
	// 2. 加密密码
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", zap.Any("error", err))
		return nil, errors.New(constant.GetErrorMessage(constant.ErrInternal))
	}

	// 3. 创建用户
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Capacity: 10 * 1024 * 1024 * 1024, // 10GB
		IsActive: true,
	}

	if err := s.userDAO.Create(user); err != nil {
		logger.Error("Failed to create user", zap.Any("error", err))
		return nil, errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("User registered successfully",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username))

	return user, nil
}

// Login 用户登录
func (s *UserServiceImpl) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	// 1. 查找用户
	user, err := s.userDAO.GetByUsernameOrEmail(req.Username)
	if err != nil {
		logger.Error("Failed to get user", zap.Any("error", err))
		return nil, errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}
	if user == nil {
		return nil, errors.New(constant.GetErrorMessage(constant.ErrUserNotExists))
	}

	// 2. 检查用户是否激活
	if !user.IsActive {
		return nil, errors.New(constant.GetErrorMessage(constant.ErrUserInactive))
	}

	// 3. 验证密码
	if !crypto.CheckPassword(req.Password, user.Password) {
		return nil, errors.New(constant.GetErrorMessage(constant.ErrPasswordWrong))
	}

	// 4. 生成JWT token
	token, expiresAt, err := s.jwtManager.Generate(user.ID, user.Username, user.Email)
	if err != nil {
		logger.Error("Failed to generate token", zap.Any("error", err))
		return nil, errors.New(constant.GetErrorMessage(constant.ErrInternal))
	}

	// 5. 更新最后登录时间
	_ = s.userDAO.UpdateLastLogin(user.ID)

	logger.Info("User logged in successfully",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username))

	return &response.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: response.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			Capacity:  user.Capacity,
			Used:      user.Used,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// GetProfile 获取用户信息
func (s *UserServiceImpl) GetProfile(userID uint) (*entity.UserResponse, error) {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		logger.Error("Failed to get user profile",
			zap.Uint("user_id", userID),
			zap.Any("error", err))
		return nil, errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}
	if user == nil {
		return nil, errors.New(constant.GetErrorMessage(constant.ErrUserNotExists))
	}

	return user.ToResponse(), nil
}

// UpdateProfile 更新个人信息
func (s *UserServiceImpl) UpdateProfile(userID uint, req *request.UpdateProfileRequest) error {
	updates := map[string]interface{}{}

	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}

	if len(updates) == 0 {
		return nil
	}

	if err := s.userDAO.UpdateFields(userID, updates); err != nil {
		logger.Error("Failed to update user profile",
			zap.Uint("user_id", userID),
			zap.Any("error", err))
		return errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("User profile updated", zap.Uint("user_id", userID))
	return nil
}

// ChangePassword 修改密码
func (s *UserServiceImpl) ChangePassword(userID uint, req *request.ChangePasswordRequest) error {
	// 1. 获取用户
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		logger.Error("Failed to get user for password change",
			zap.Uint("user_id", userID),
			zap.Any("error", err))
		return errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}
	if user == nil {
		return errors.New(constant.GetErrorMessage(constant.ErrUserNotExists))
	}

	// 2. 验证旧密码
	if !crypto.CheckPassword(req.OldPassword, user.Password) {
		return errors.New(constant.GetErrorMessage(constant.ErrPasswordWrong))
	}

	// 3. 加密新密码
	hashedPassword, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		logger.Error("Failed to hash new password", zap.Any("error", err))
		return errors.New(constant.GetErrorMessage(constant.ErrInternal))
	}

	// 4. 更新密码
	if err := s.userDAO.UpdatePassword(userID, hashedPassword); err != nil {
		logger.Error("Failed to update password",
			zap.Uint("user_id", userID),
			zap.Any("error", err))
		return errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("User password changed", zap.Uint("user_id", userID))
	return nil
}

// GetStorageInfo 获取存储空间信息
func (s *UserServiceImpl) GetStorageInfo(userID uint) (*response.StorageInfo, error) {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		logger.Error("Failed to get user for storage info",
			zap.Uint("user_id", userID),
			zap.Any("error", err))
		return nil, errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}
	if user == nil {
		return nil, errors.New(constant.GetErrorMessage(constant.ErrUserNotExists))
	}

	available := user.Capacity - user.Used
	var usageRate float64
	if user.Capacity > 0 {
		usageRate = float64(user.Used) / float64(user.Capacity) * 100
	}

	// TODO: 获取文件数和文件夹数（后续实现）
	fileCount := int64(0)
	folderCount := int64(0)

	return &response.StorageInfo{
		Capacity:    user.Capacity,
		Used:        user.Used,
		Available:   available,
		UsageRate:   usageRate,
		FileCount:   fileCount,
		FolderCount: folderCount,
	}, nil
}

// CheckStorage 检查存储空间
func (s *UserServiceImpl) CheckStorage(userID uint, fileSize int64) (bool, error) {
	ok, err := s.userDAO.CheckStorage(userID, fileSize)
	if err != nil {
		logger.Error("Failed to check storage",
			zap.Uint("user_id", userID),
			zap.Int64("file_size", fileSize),
			zap.Any("error", err))
		return false, errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}
	return ok, nil
}

// AddUsedSpace 增加已用空间
func (s *UserServiceImpl) AddUsedSpace(userID uint, size int64) error {
	if err := s.userDAO.UpdateStorageUsed(userID, size); err != nil {
		logger.Error("Failed to add used space",
			zap.Uint("user_id", userID),
			zap.Int64("size", size),
			zap.Any("error", err))
		return errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}
	return nil
}

// ReduceUsedSpace 减少已用空间
func (s *UserServiceImpl) ReduceUsedSpace(userID uint, size int64) error {
	if err := s.userDAO.UpdateStorageUsed(userID, -size); err != nil {
		logger.Error("Failed to reduce used space",
			zap.Uint("user_id", userID),
			zap.Int64("size", size),
			zap.Any("error", err))
		return errors.New(constant.GetErrorMessage(constant.ErrDatabase))
	}
	return nil
}
