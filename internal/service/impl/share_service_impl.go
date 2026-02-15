// internal/service/impl/share_service_impl.go
package impl

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"cloud-disk/internal/constant"
	"cloud-disk/internal/dao"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
	"cloud-disk/internal/model/entity"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/internal/utils/crypto"
	"cloud-disk/pkg/logger"
	"cloud-disk/pkg/storage"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ShareServiceImpl struct {
	shareDAO *dao.ShareDAO
	fileDAO  *dao.FileDAO
	userDAO  *dao.UserDAO
	storage  storage.Storage
	db       *gorm.DB
	baseURL  string
}

func NewShareService(db *gorm.DB, shareDAO *dao.ShareDAO, fileDAO *dao.FileDAO, userDAO *dao.UserDAO, storage storage.Storage) interfaces.ShareService {
	return &ShareServiceImpl{
		shareDAO: shareDAO,
		fileDAO:  fileDAO,
		userDAO:  userDAO,
		storage:  storage,
		db:       db,
		baseURL:  "http://localhost:8080/s/", // 分享链接基础URL
	}
}

// generateToken 生成唯一分享令牌
func (s *ShareServiceImpl) generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Create 创建分享
func (s *ShareServiceImpl) Create(ctx context.Context, userID uint, req *request.ShareCreateRequest) (*response.ShareResponse, error) {
	// 1. 检查文件是否存在且属于该用户
	file, err := s.fileDAO.GetByID(req.FileID)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if file == nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}
	if file.UserID != userID {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
	}

	// 2. 检查是否已存在有效分享
	existingShare, err := s.shareDAO.GetByFileID(req.FileID)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if existingShare != nil && existingShare.Status == 1 {
		return nil, fmt.Errorf("文件已存在有效分享")
	}

	// 3. 生成唯一token
	token, err := s.generateToken()
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrInternal))
	}

	// 4. 处理过期时间
	var expireTime *time.Time
	if req.ExpireDays > 0 {
		et := time.Now().AddDate(0, 0, req.ExpireDays)
		expireTime = &et
	}

	// 5. 处理密码
	var hashedPassword string
	if req.Password != "" {
		hashedPassword, err = crypto.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrInternal))
		}
	}

	// 6. 创建分享记录
	share := &entity.Share{
		Token:        token,
		FileID:       req.FileID,
		UserID:       userID,
		Password:     hashedPassword,
		ExpireTime:   expireTime,
		MaxDownloads: req.MaxDownloads,
		Status:       1,
	}

	if err := s.shareDAO.Create(share); err != nil {
		logger.Error("Failed to create share", zap.Error(err))
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 7. 构建响应
	shareLink := fmt.Sprintf("%s%s", s.baseURL, token)
	shareResp := &response.ShareResponse{
		ID:            share.ID,
		Token:         share.Token,
		FileID:        share.FileID,
		Filename:      file.Filename,
		FileSize:      file.Size,
		FileSizeText:  entity.FormatFileSize(file.Size),
		ExpireTime:    share.ExpireTime,
		MaxDownloads:  share.MaxDownloads,
		DownloadCount: share.DownloadCount,
		Status:        share.Status,
		ShareLink:     shareLink,
		CreatedAt:     share.CreatedAt,
	}

	logger.Info("Share created successfully",
		zap.Uint("share_id", share.ID),
		zap.String("token", token),
		zap.Uint("user_id", userID))

	return shareResp, nil
}

// Cancel 取消分享
func (s *ShareServiceImpl) Cancel(ctx context.Context, userID uint, shareID uint) error {
	share, err := s.shareDAO.GetByID(shareID)
	if err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if share == nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrShareNotFound))
	}
	if share.UserID != userID {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
	}

	// 更新状态为已取消
	if err := s.shareDAO.UpdateFields(shareID, map[string]interface{}{
		"status": 2,
	}); err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("Share cancelled", zap.Uint("share_id", shareID))
	return nil
}

// Update 更新分享
func (s *ShareServiceImpl) Update(ctx context.Context, userID uint, shareID uint, req *request.ShareUpdateRequest) error {
	share, err := s.shareDAO.GetByID(shareID)
	if err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if share == nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrShareNotFound))
	}
	if share.UserID != userID {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
	}

	updates := map[string]interface{}{}

	// 更新密码
	if req.Password != nil {
		if *req.Password == "" {
			updates["password"] = ""
		} else {
			hashedPassword, err := crypto.HashPassword(*req.Password)
			if err != nil {
				return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrInternal))
			}
			updates["password"] = hashedPassword
		}
	}

	// 更新过期时间
	if req.ExpireDays != nil {
		if *req.ExpireDays <= 0 {
			updates["expire_time"] = nil
		} else {
			et := time.Now().AddDate(0, 0, *req.ExpireDays)
			updates["expire_time"] = et
		}
	}

	// 更新下载次数限制
	if req.MaxDownloads != nil {
		updates["max_downloads"] = *req.MaxDownloads
	}

	// 更新状态
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := s.shareDAO.UpdateFields(shareID, updates); err != nil {
			return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
		}
	}

	logger.Info("Share updated", zap.Uint("share_id", shareID))
	return nil
}

// GetDetail 获取分享详情
func (s *ShareServiceImpl) GetDetail(ctx context.Context, userID uint, shareID uint) (*response.ShareResponse, error) {
	share, err := s.shareDAO.GetByID(shareID)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if share == nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrShareNotFound))
	}
	if share.UserID != userID {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
	}

	// 获取文件信息
	file, err := s.fileDAO.GetByID(share.FileID)
	if err != nil || file == nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}

	shareLink := fmt.Sprintf("%s%s", s.baseURL, share.Token)
	shareResp := &response.ShareResponse{
		ID:            share.ID,
		Token:         share.Token,
		FileID:        share.FileID,
		Filename:      file.Filename,
		FileSize:      file.Size,
		FileSizeText:  entity.FormatFileSize(file.Size),
		ExpireTime:    share.ExpireTime,
		MaxDownloads:  share.MaxDownloads,
		DownloadCount: share.DownloadCount,
		Status:        share.Status,
		ShareLink:     shareLink,
		CreatedAt:     share.CreatedAt,
	}

	return shareResp, nil
}

// GetList 获取用户分享列表
func (s *ShareServiceImpl) GetList(ctx context.Context, userID uint, req *request.ShareListRequest) (*response.ShareListResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	shares, total, err := s.shareDAO.ListByUserID(userID, req.Status, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	shareResponses := make([]*response.ShareResponse, 0, len(shares))

	for _, share := range shares {
		// 获取文件信息
		file, err := s.fileDAO.GetByID(share.FileID)
		if err != nil || file == nil {
			continue
		}

		shareLink := fmt.Sprintf("%s%s", s.baseURL, share.Token)
		shareResp := &response.ShareResponse{
			ID:            share.ID,
			Token:         share.Token,
			FileID:        share.FileID,
			Filename:      file.Filename,
			FileSize:      file.Size,
			FileSizeText:  entity.FormatFileSize(file.Size),
			ExpireTime:    share.ExpireTime,
			MaxDownloads:  share.MaxDownloads,
			DownloadCount: share.DownloadCount,
			Status:        share.Status,
			ShareLink:     shareLink,
			CreatedAt:     share.CreatedAt,
		}
		shareResponses = append(shareResponses, shareResp)
	}

	return &response.ShareListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Shares:   shareResponses,
	}, nil
}

// AccessByToken 通过token访问分享
func (s *ShareServiceImpl) AccessByToken(ctx context.Context, token, password string) (*response.ShareDetailResponse, error) {
	// 1. 获取分享信息
	share, err := s.shareDAO.GetByToken(token)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if share == nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrShareNotFound))
	}

	// 2. 验证分享状态
	if share.Status != 1 {
		return nil, fmt.Errorf("分享已取消")
	}

	if share.IsExpired() {
		return nil, fmt.Errorf("分享已过期")
	}

	if share.IsDownloadLimitReached() {
		return nil, fmt.Errorf("分享下载次数已达上限")
	}

	// 3. 验证密码
	needPassword := share.Password != ""
	if needPassword {
		if password == "" {
			return &response.ShareDetailResponse{
				Token:        share.Token,
				NeedPassword: true,
			}, nil
		}
		if !crypto.CheckPassword(password, share.Password) {
			return nil, fmt.Errorf("密码错误")
		}
	}

	// 4. 获取文件信息
	file, err := s.fileDAO.GetByID(share.FileID)
	if err != nil || file == nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}

	// 5. 计算剩余下载次数
	downloadLeft := -1
	if share.MaxDownloads > 0 {
		downloadLeft = share.MaxDownloads - share.DownloadCount
	}

	// 6. 构建文件响应
	fileResp := &response.FileResponse{
		ID:        file.ID,
		Filename:  file.Filename,
		Path:      file.Path,
		Size:      file.Size,
		SizeText:  entity.FormatFileSize(file.Size),
		MimeType:  file.MimeType,
		Extension: file.Extension,
		IsDir:     file.IsDir,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}

	return &response.ShareDetailResponse{
		Token:        share.Token,
		File:         fileResp,
		ExpireTime:   share.ExpireTime,
		DownloadLeft: downloadLeft,
		NeedPassword: false,
		CreatedAt:    share.CreatedAt,
	}, nil
}

// DownloadShare 下载分享文件
func (s *ShareServiceImpl) DownloadShare(ctx context.Context, token, password string) (string, error) {
	// 1. 获取分享信息
	share, err := s.shareDAO.GetByToken(token)
	if err != nil {
		return "", fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if share == nil {
		return "", fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrShareNotFound))
	}

	// 2. 验证分享状态
	if share.Status != 1 {
		return "", fmt.Errorf("分享已取消")
	}

	if share.IsExpired() {
		return "", fmt.Errorf("分享已过期")
	}

	if share.IsDownloadLimitReached() {
		return "", fmt.Errorf("分享下载次数已达上限")
	}

	// 3. 验证密码
	if share.Password != "" {
		if password == "" {
			return "", fmt.Errorf("需要密码")
		}
		if !crypto.CheckPassword(password, share.Password) {
			return "", fmt.Errorf("密码错误")
		}
	}

	// 4. 获取文件信息
	file, err := s.fileDAO.GetByID(share.FileID)
	if err != nil || file == nil {
		return "", fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}

	// 5. 增加下载次数
	if err := s.shareDAO.IncrementDownloadCount(token); err != nil {
		logger.Error("Failed to increment download count", zap.Error(err))
	}

	// 6. 获取文件下载URL
	url, err := s.storage.GetURL(ctx, file.Path)
	if err != nil {
		return "", fmt.Errorf("获取文件URL失败")
	}

	return url, nil
}

// ValidateAccess 验证分享访问权限
func (s *ShareServiceImpl) ValidateAccess(ctx context.Context, token, password string) (string, error) {
	share, err := s.shareDAO.GetByToken(token)
	if err != nil {
		return "", fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if share == nil {
		return "", fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrShareNotFound))
	}

	// 验证密码
	if share.Password != "" {
		if password == "" {
			return "", fmt.Errorf("需要密码")
		}
		if !crypto.CheckPassword(password, share.Password) {
			return "", fmt.Errorf("密码错误")
		}
	}

	return share.FileID, nil
}
