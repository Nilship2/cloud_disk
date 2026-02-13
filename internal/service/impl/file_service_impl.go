// internal/service/impl/file_service_impl.go
package impl

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"cloud-disk/internal/constant"
	"cloud-disk/internal/dao"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
	"cloud-disk/internal/model/entity"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/pkg/logger"
	storageInterface "cloud-disk/pkg/storage" // 使用别名

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FileServiceImpl struct {
	fileDAO *dao.FileDAO
	userDAO *dao.UserDAO
	storage storageInterface.Storage // 使用接口类型
	db      *gorm.DB
}

func NewFileService(db *gorm.DB, fileDAO *dao.FileDAO, userDAO *dao.UserDAO, storage storageInterface.Storage) interfaces.FileService {
	return &FileServiceImpl{
		fileDAO: fileDAO,
		userDAO: userDAO,
		storage: storage,
		db:      db,
	}
}

// Upload 上传文件
func (s *FileServiceImpl) Upload(ctx context.Context, userID uint, file *multipart.FileHeader, parentID string) (*response.FileResponse, error) {
	// 1. 检查存储空间
	ok, err := s.userDAO.CheckStorage(userID, file.Size)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if !ok {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrStorageFull))
	}

	// 2. 生成文件ID
	fileID := uuid.New().String()

	// 3. 计算文件哈希（用于秒传）
	hash, err := s.calculateFileHash(file)
	if err != nil {
		logger.Error("Failed to calculate file hash", zap.Error(err))
	}

	// 4. 上传文件到存储
	path, size, err := s.storage.Upload(ctx, file, userID, fileID)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileUploadFailed))
	}

	// 5. 创建文件记录
	fileEntity := &entity.File{
		ID:        fileID,
		UserID:    userID,
		Filename:  file.Filename,
		Path:      path,
		Size:      size,
		Hash:      hash,
		MimeType:  file.Header.Get("Content-Type"),
		Extension: strings.TrimPrefix(filepath.Ext(file.Filename), "."),
		IsDir:     false,
		Status:    1,
	}

	// 设置父目录
	if parentID != "" {
		parent, err := s.fileDAO.GetByID(parentID)
		if err != nil || parent == nil {
			return nil, fmt.Errorf("parent directory not found")
		}
		if !parent.IsDir {
			return nil, fmt.Errorf("parent is not a directory")
		}
		fileEntity.ParentID = &parentID
	}

	// 6. 开启数据库事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// 保存文件记录
	if err := s.fileDAO.Create(fileEntity); err != nil {
		tx.Rollback()
		// 上传失败，删除已存储的文件
		s.storage.Delete(ctx, path)
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 更新用户已用空间
	if err := s.userDAO.UpdateStorageUsed(userID, size); err != nil {
		tx.Rollback()
		s.storage.Delete(ctx, path)
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		s.storage.Delete(ctx, path)
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	url, _ := s.storage.GetURL(ctx, path)
	fileResponse := &response.FileResponse{
		ID:        fileEntity.ID,
		Filename:  fileEntity.Filename,
		Path:      fileEntity.Path,
		Size:      fileEntity.Size,
		SizeText:  entity.FormatFileSize(fileEntity.Size),
		MimeType:  fileEntity.MimeType,
		Extension: fileEntity.Extension,
		IsDir:     fileEntity.IsDir,
		ParentID:  fileEntity.ParentID,
		CreatedAt: fileEntity.CreatedAt,
		UpdatedAt: fileEntity.UpdatedAt,
		URL:       url,
	}

	logger.Info("File uploaded successfully",
		zap.String("file_id", fileEntity.ID),
		zap.Uint("user_id", userID),
		zap.Int64("size", fileEntity.Size))

	return fileResponse, nil
}

// InstantUpload 秒传文件
func (s *FileServiceImpl) InstantUpload(ctx context.Context, userID uint, hash, filename string, parentID string) (*response.FileResponse, error) {
	// 1. 查找是否存在相同哈希的文件
	existingFile, err := s.fileDAO.GetByHash(hash)
	if err != nil || existingFile == nil {
		return nil, fmt.Errorf("file not found for instant upload")
	}

	// 2. 检查存储空间
	ok, err := s.userDAO.CheckStorage(userID, existingFile.Size)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if !ok {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrStorageFull))
	}

	// 3. 生成新的文件ID
	fileID := uuid.New().String()

	// 4. 创建文件记录（引用同一物理文件）
	fileEntity := &entity.File{
		ID:        fileID,
		UserID:    userID,
		Filename:  filename,
		Path:      existingFile.Path, // 使用相同路径
		Size:      existingFile.Size,
		Hash:      hash,
		MimeType:  existingFile.MimeType,
		Extension: existingFile.Extension,
		IsDir:     false,
		Status:    1,
	}

	if parentID != "" {
		fileEntity.ParentID = &parentID
	}

	// 5. 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.fileDAO.Create(fileEntity); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	if err := s.userDAO.UpdateStorageUsed(userID, existingFile.Size); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	url, _ := s.storage.GetURL(ctx, existingFile.Path)
	fileResponse := &response.FileResponse{
		ID:        fileEntity.ID,
		Filename:  fileEntity.Filename,
		Path:      fileEntity.Path,
		Size:      fileEntity.Size,
		SizeText:  entity.FormatFileSize(fileEntity.Size),
		MimeType:  fileEntity.MimeType,
		Extension: fileEntity.Extension,
		IsDir:     fileEntity.IsDir,
		ParentID:  fileEntity.ParentID,
		CreatedAt: fileEntity.CreatedAt,
		UpdatedAt: fileEntity.UpdatedAt,
		URL:       url,
	}

	logger.Info("File instant uploaded successfully",
		zap.String("file_id", fileEntity.ID),
		zap.Uint("user_id", userID),
		zap.String("hash", hash))

	return fileResponse, nil
}

// Download 下载文件
func (s *FileServiceImpl) Download(ctx context.Context, userID uint, fileID string) (string, error) {
	// 1. 获取文件信息
	file, err := s.fileDAO.GetByID(fileID)
	if err != nil {
		return "", fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if file == nil {
		return "", fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}

	// 2. 权限检查：只能下载自己的文件
	if file.UserID != userID {
		return "", fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
	}

	// 3. 获取文件访问URL
	url, err := s.storage.GetURL(ctx, file.Path)
	if err != nil {
		return "", fmt.Errorf("failed to get file URL")
	}

	return url, nil
}

// GetList 获取文件列表
func (s *FileServiceImpl) GetList(ctx context.Context, userID uint, req *request.FileListRequest) (*response.FileListResponse, error) {
	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	order := req.Order
	if order == "" {
		order = "desc"
	}

	var parentID *string
	if req.ParentID != "" {
		parentID = &req.ParentID
	}

	// 查询文件列表
	files, total, err := s.fileDAO.ListByUserID(userID, parentID, page, pageSize, orderBy, order)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 转换为响应DTO - 使用 response.FileResponse
	fileResponses := make([]*response.FileResponse, 0, len(files))
	for _, file := range files {
		resp := &response.FileResponse{
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
		// 获取访问URL
		url, _ := s.storage.GetURL(ctx, file.Path)
		resp.URL = url
		fileResponses = append(fileResponses, resp)
	}

	return &response.FileListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Files:    fileResponses,
	}, nil
}

// GetDetail 获取文件详情
func (s *FileServiceImpl) GetDetail(ctx context.Context, userID uint, fileID string) (*response.FileResponse, error) {
	file, err := s.fileDAO.GetByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if file == nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}

	if file.UserID != userID {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
	}

	url, _ := s.storage.GetURL(ctx, file.Path)
	fileResponse := &response.FileResponse{
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
		URL:       url,
	}

	return fileResponse, nil
}

// Delete 删除文件（移入回收站）
func (s *FileServiceImpl) Delete(ctx context.Context, userID uint, fileID string) error {
	file, err := s.fileDAO.GetByID(fileID)
	if err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if file == nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}

	if file.UserID != userID {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
	}

	// 软删除
	if err := s.fileDAO.Delete(fileID); err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 减少用户已用空间
	if err := s.userDAO.UpdateStorageUsed(userID, -file.Size); err != nil {
		logger.Error("Failed to update storage used", zap.Error(err))
	}

	logger.Info("File moved to trash",
		zap.String("file_id", fileID),
		zap.Uint("user_id", userID))

	return nil
}

// BatchDelete 批量删除
func (s *FileServiceImpl) BatchDelete(ctx context.Context, userID uint, fileIDs []string) error {
	var totalSize int64 = 0

	// 验证所有文件都属于该用户并计算总大小
	for _, fileID := range fileIDs {
		file, err := s.fileDAO.GetByID(fileID)
		if err != nil {
			return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
		}
		if file == nil {
			return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
		}
		if file.UserID != userID {
			return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
		}
		totalSize += file.Size
	}

	// 批量软删除
	if err := s.fileDAO.BatchDelete(fileIDs); err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	// 减少用户已用空间
	if err := s.userDAO.UpdateStorageUsed(userID, -totalSize); err != nil {
		logger.Error("Failed to update storage used", zap.Error(err))
	}

	logger.Info("Files batch moved to trash",
		zap.Int("count", len(fileIDs)),
		zap.Uint("user_id", userID))

	return nil
}

// CreateFolder 创建文件夹
func (s *FileServiceImpl) CreateFolder(ctx context.Context, userID uint, name, parentID string) (*response.FileResponse, error) {
	// 检查同一目录下是否已存在同名文件夹
	var parentIDPtr *string
	if parentID != "" {
		parentIDPtr = &parentID
	}

	files, _, err := s.fileDAO.ListByUserID(userID, parentIDPtr, 1, 1, "filename", "asc")
	if err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	for _, f := range files {
		if f.IsDir && f.Filename == name {
			return nil, fmt.Errorf("folder with same name already exists")
		}
	}

	// 创建文件夹记录
	folderID := uuid.New().String()
	folder := &entity.File{
		ID:       folderID,
		UserID:   userID,
		Filename: name,
		Path:     fmt.Sprintf("user_%d/%s", userID, folderID),
		Size:     0,
		IsDir:    true,
		Status:   1,
		ParentID: parentIDPtr,
	}

	if err := s.fileDAO.Create(folder); err != nil {
		return nil, fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	folderResponse := &response.FileResponse{
		ID:        folder.ID,
		Filename:  folder.Filename,
		Path:      folder.Path,
		Size:      folder.Size,
		SizeText:  entity.FormatFileSize(folder.Size),
		IsDir:     folder.IsDir,
		ParentID:  folder.ParentID,
		CreatedAt: folder.CreatedAt,
		UpdatedAt: folder.UpdatedAt,
	}

	logger.Info("Folder created successfully",
		zap.String("folder_id", folder.ID),
		zap.Uint("user_id", userID),
		zap.String("name", name))

	return folderResponse, nil
}

// Rename 重命名
func (s *FileServiceImpl) Rename(ctx context.Context, userID uint, fileID, newName string) error {
	file, err := s.fileDAO.GetByID(fileID)
	if err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if file == nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}
	if file.UserID != userID {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
	}

	// 更新文件名
	if err := s.fileDAO.UpdateFields(fileID, map[string]interface{}{
		"filename": newName,
	}); err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("File renamed successfully",
		zap.String("file_id", fileID),
		zap.String("old_name", file.Filename),
		zap.String("new_name", newName))

	return nil
}

// Move 移动文件
func (s *FileServiceImpl) Move(ctx context.Context, userID uint, fileID, targetParentID string) error {
	file, err := s.fileDAO.GetByID(fileID)
	if err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}
	if file == nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrFileNotFound))
	}
	if file.UserID != userID {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrPermissionDenied))
	}

	// 检查目标目录是否存在
	if targetParentID != "" {
		target, err := s.fileDAO.GetByID(targetParentID)
		if err != nil || target == nil {
			return fmt.Errorf("target directory not found")
		}
		if !target.IsDir {
			return fmt.Errorf("target is not a directory")
		}
	}

	// 更新父目录
	var parentID *string
	if targetParentID != "" {
		parentID = &targetParentID
	}

	if err := s.fileDAO.UpdateFields(fileID, map[string]interface{}{
		"parent_id": parentID,
	}); err != nil {
		return fmt.Errorf("%s", constant.GetErrorMessage(constant.ErrDatabase))
	}

	logger.Info("File moved successfully",
		zap.String("file_id", fileID),
		zap.String("target_parent_id", targetParentID))

	return nil
}

// CheckExists 检查文件是否存在（用于秒传）
func (s *FileServiceImpl) CheckExists(ctx context.Context, userID uint, hash string) (bool, *response.FileResponse) {
	file, err := s.fileDAO.GetByHash(hash)
	if err != nil || file == nil {
		return false, nil
	}

	fileResponse := &response.FileResponse{
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

	return true, fileResponse
}

// calculateFileHash 计算文件哈希值
func (s *FileServiceImpl) calculateFileHash(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
