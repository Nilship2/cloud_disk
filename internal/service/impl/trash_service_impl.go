// internal/service/impl/trash_service_impl.go
package impl

import (
	"context"
	"fmt"
	"time"

	"cloud-disk/internal/constant"
	"cloud-disk/internal/dao"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
	"cloud-disk/internal/model/entity"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/pkg/logger"
	"cloud-disk/pkg/storage"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TrashServiceImpl struct {
	fileDAO *dao.FileDAO
	userDAO *dao.UserDAO
	storage storage.Storage
	db      *gorm.DB
}

func NewTrashService(db *gorm.DB, fileDAO *dao.FileDAO, userDAO *dao.UserDAO, storage storage.Storage) interfaces.TrashService {
	return &TrashServiceImpl{
		fileDAO: fileDAO,
		userDAO: userDAO,
		storage: storage,
		db:      db,
	}
}

// GetList 获取回收站列表
func (s *TrashServiceImpl) GetList(ctx context.Context, userID uint, req *request.TrashListRequest) (*response.TrashListResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// 获取回收站列表
	files, total, err := s.fileDAO.GetTrashList(userID, page, pageSize)
	if err != nil {
		logger.Error("Failed to get trash list", zap.Error(err))
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 转换为响应DTO
	items := make([]*response.TrashItemResponse, 0, len(files))
	for _, file := range files {
		// 计算剩余过期天数
		//expireDays := 30
		var expireIn int
		if file.DeletedAt.Valid {
			expireTime := file.DeletedAt.Time.AddDate(0, 0, 30)
			daysLeft := int(time.Until(expireTime).Hours() / 24)
			if daysLeft > 0 {
				expireIn = daysLeft
			} else {
				expireIn = 0
			}
		}

		item := &response.TrashItemResponse{
			ID:        file.ID,
			Filename:  file.Filename,
			Size:      file.Size,
			SizeText:  entity.FormatFileSize(file.Size),
			MimeType:  file.MimeType,
			Extension: file.Extension,
			IsDir:     file.IsDir,
			DeletedAt: file.DeletedAt.Time,
			ExpireIn:  expireIn,
		}
		items = append(items, item)
	}

	return &response.TrashListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}, nil
}

// Restore 恢复文件
func (s *TrashServiceImpl) Restore(ctx context.Context, userID uint, fileID string) error {
	// 获取回收站文件
	file, err := s.fileDAO.GetTrashItem(userID, fileID)
	if err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if file == nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 恢复文件（软删除恢复）
	if err := s.fileDAO.RestoreFromTrash(fileID); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 恢复用户已用空间
	if err := s.userDAO.UpdateStorageUsed(userID, file.Size); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("File restored from trash",
		zap.String("file_id", fileID),
		zap.Uint("user_id", userID))

	return nil
}

// BatchRestore 批量恢复
func (s *TrashServiceImpl) BatchRestore(ctx context.Context, userID uint, req *request.TrashBatchRequest) error {
	var totalSize int64 = 0
	var validFileIDs []string

	// 验证所有文件并计算总大小
	for _, fileID := range req.FileIDs {
		file, err := s.fileDAO.GetTrashItem(userID, fileID)
		if err != nil {
			return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
		}
		if file != nil {
			totalSize += file.Size
			validFileIDs = append(validFileIDs, fileID)
		}
	}

	if len(validFileIDs) == 0 {
		return fmt.Errorf("没有可恢复的文件")
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 批量恢复
	if err := s.fileDAO.BatchRestoreFromTrash(validFileIDs); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 恢复用户已用空间
	if err := s.userDAO.UpdateStorageUsed(userID, totalSize); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("Batch restored from trash",
		zap.Int("count", len(validFileIDs)),
		zap.Uint("user_id", userID))

	return nil
}

// Delete 彻底删除
func (s *TrashServiceImpl) Delete(ctx context.Context, userID uint, fileID string) error {
	// 获取回收站文件
	file, err := s.fileDAO.GetTrashItem(userID, fileID)
	if err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if file == nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除物理文件
	if err := s.storage.Delete(ctx, file.Path); err != nil {
		logger.Warn("Failed to delete physical file", zap.Error(err))
		// 继续执行，即使物理文件删除失败也尝试删除数据库记录
	}

	// 彻底删除数据库记录
	if err := s.fileDAO.PermanentlyDelete(fileID); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 注意：不减少用户空间，因为已经在移入回收站时减少了

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("File permanently deleted",
		zap.String("file_id", fileID),
		zap.Uint("user_id", userID))

	return nil
}

// BatchDelete 批量彻底删除
func (s *TrashServiceImpl) BatchDelete(ctx context.Context, userID uint, req *request.TrashBatchRequest) error {
	var validFileIDs []string

	// 验证所有文件
	for _, fileID := range req.FileIDs {
		file, err := s.fileDAO.GetTrashItem(userID, fileID)
		if err != nil {
			return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
		}
		if file != nil {
			validFileIDs = append(validFileIDs, fileID)
			// 尝试删除物理文件（不阻塞主流程）
			go func(path string) {
				if err := s.storage.Delete(ctx, path); err != nil {
					logger.Warn("Failed to delete physical file", zap.Error(err))
				}
			}(file.Path)
		}
	}

	if len(validFileIDs) == 0 {
		return fmt.Errorf("没有可删除的文件")
	}

	// 批量彻底删除数据库记录
	if err := s.fileDAO.BatchPermanentlyDelete(validFileIDs); err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("Batch permanently deleted",
		zap.Int("count", len(validFileIDs)),
		zap.Uint("user_id", userID))

	return nil
}

// CleanAll 清空回收站
func (s *TrashServiceImpl) CleanAll(ctx context.Context, userID uint) (*response.TrashCleanResponse, error) {
	// 获取所有回收站文件
	files, _, err := s.fileDAO.GetTrashList(userID, 1, 1000) // 假设最多1000个
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	if len(files) == 0 {
		return &response.TrashCleanResponse{
			CleanedCount:   0,
			FreedSpace:     0,
			FreedSpaceText: "0 B",
		}, nil
	}

	// 计算总大小
	var totalSize int64
	var fileIDs []string
	for _, file := range files {
		totalSize += file.Size
		fileIDs = append(fileIDs, file.ID)
		// 异步删除物理文件
		go func(path string) {
			if err := s.storage.Delete(ctx, path); err != nil {
				logger.Warn("Failed to delete physical file", zap.Error(err))
			}
		}(file.Path)
	}

	// 批量彻底删除
	if err := s.fileDAO.BatchPermanentlyDelete(fileIDs); err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("Trash cleaned",
		zap.Int("count", len(files)),
		zap.Int64("size", totalSize),
		zap.Uint("user_id", userID))

	return &response.TrashCleanResponse{
		CleanedCount:   int64(len(files)),
		FreedSpace:     totalSize,
		FreedSpaceText: entity.FormatFileSize(totalSize),
	}, nil
}

// GetStats 获取回收站统计信息
func (s *TrashServiceImpl) GetStats(ctx context.Context, userID uint) (int64, int64, error) {
	count, err := s.fileDAO.CountTrashByUser(userID)
	if err != nil {
		return 0, 0, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	size, err := s.fileDAO.SumTrashSizeByUser(userID)
	if err != nil {
		return 0, 0, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	return count, size, nil
}
