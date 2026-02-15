// internal/service/impl/favorite_service_impl.go
package impl

import (
	"context"
	"fmt"

	"cloud-disk/internal/constant"
	"cloud-disk/internal/dao"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
	"cloud-disk/internal/model/entity"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FavoriteServiceImpl struct {
	favoriteDAO *dao.FavoriteDAO
	fileDAO     *dao.FileDAO
	db          *gorm.DB
}

func NewFavoriteService(db *gorm.DB, favoriteDAO *dao.FavoriteDAO, fileDAO *dao.FileDAO) interfaces.FavoriteService {
	return &FavoriteServiceImpl{
		favoriteDAO: favoriteDAO,
		fileDAO:     fileDAO,
		db:          db,
	}
}

// Add 添加收藏
func (s *FavoriteServiceImpl) Add(ctx context.Context, userID uint, req *request.FavoriteCreateRequest) (*response.FavoriteResponse, error) {
	// 1. 检查文件是否存在
	file, err := s.fileDAO.GetByID(req.FileID)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if file == nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}

	// 2. 检查是否已收藏
	exists, err := s.favoriteDAO.IsFavorite(userID, req.FileID)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if exists {
		return nil, fmt.Errorf("文件已收藏")
	}

	// 3. 创建收藏记录
	favorite := &entity.Favorite{
		UserID: userID,
		FileID: req.FileID,
	}

	if err := s.favoriteDAO.Create(favorite); err != nil {
		logger.Error("Failed to create favorite", zap.Error(err))
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 4. 构建响应
	fileResp := &response.FileResponse{
		ID:        file.ID,
		Filename:  file.Filename,
		Path:      file.Path,
		Size:      file.Size,
		SizeText:  entity.FormatFileSize(file.Size),
		MimeType:  file.MimeType,
		Extension: file.Extension,
		IsDir:     file.IsDir,
		ParentID:  file.ParentID,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}

	favoriteResp := &response.FavoriteResponse{
		ID:        favorite.ID,
		File:      fileResp,
		CreatedAt: favorite.CreatedAt,
	}

	logger.Info("Favorite added",
		zap.Uint("user_id", userID),
		zap.String("file_id", req.FileID))

	return favoriteResp, nil
}

// Remove 取消收藏
func (s *FavoriteServiceImpl) Remove(ctx context.Context, userID uint, fileID string) error {
	// 检查是否存在
	favorite, err := s.favoriteDAO.GetByUserAndFile(userID, fileID)
	if err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if favorite == nil {
		return fmt.Errorf("收藏不存在")
	}

	// 删除收藏
	if err := s.favoriteDAO.Delete(userID, fileID); err != nil {
		logger.Error("Failed to remove favorite", zap.Error(err))
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("Favorite removed",
		zap.Uint("user_id", userID),
		zap.String("file_id", fileID))

	return nil
}

// GetList 获取收藏列表
func (s *FavoriteServiceImpl) GetList(ctx context.Context, userID uint, req *request.FavoriteListRequest) (*response.FavoriteListResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// 获取收藏列表（包含文件信息）
	favorites, total, err := s.favoriteDAO.ListByUserID(userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	favoriteResponses := make([]*response.FavoriteResponse, 0, len(favorites))

	for _, favorite := range favorites {
		// 获取文件信息
		file, err := s.fileDAO.GetByID(favorite.FileID)
		if err != nil || file == nil {
			continue
		}

		fileResp := &response.FileResponse{
			ID:        file.ID,
			Filename:  file.Filename,
			Path:      file.Path,
			Size:      file.Size,
			SizeText:  entity.FormatFileSize(file.Size),
			MimeType:  file.MimeType,
			Extension: file.Extension,
			IsDir:     file.IsDir,
			ParentID:  file.ParentID,
			CreatedAt: file.CreatedAt,
			UpdatedAt: file.UpdatedAt,
		}

		favoriteResp := &response.FavoriteResponse{
			ID:        favorite.ID,
			File:      fileResp,
			CreatedAt: favorite.CreatedAt,
		}

		favoriteResponses = append(favoriteResponses, favoriteResp)
	}

	return &response.FavoriteListResponse{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		Favorites: favoriteResponses,
	}, nil
}

// IsFavorite 检查是否已收藏
func (s *FavoriteServiceImpl) IsFavorite(ctx context.Context, userID uint, fileID string) (bool, error) {
	exists, err := s.favoriteDAO.IsFavorite(userID, fileID)
	if err != nil {
		return false, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	return exists, nil
}
