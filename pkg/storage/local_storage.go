// pkg/storage/local_storage.go
package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud-disk/pkg/logger"

	"go.uber.org/zap"
)

type LocalStorage struct {
	basePath  string
	tempPath  string
	maxSizeMB int64
	baseURL   string
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(basePath, tempPath string, maxSizeMB int64) *LocalStorage {
	// 确保目录存在
	os.MkdirAll(basePath, 0755)
	os.MkdirAll(tempPath, 0755)

	return &LocalStorage{
		basePath:  basePath,
		tempPath:  tempPath,
		maxSizeMB: maxSizeMB,
		baseURL:   "/files", // 静态文件服务路径
	}
}

// Upload 上传文件
func (s *LocalStorage) Upload(ctx context.Context, file *multipart.FileHeader, userID uint, fileID string) (string, int64, error) {
	// 1. 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// 2. 生成存储路径：basePath/user_id/file_id/filename
	userDir := filepath.Join(s.basePath, fmt.Sprintf("user_%d", userID))
	fileDir := filepath.Join(userDir, fileID)

	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return "", 0, fmt.Errorf("failed to create directory: %w", err)
	}

	// 3. 安全处理文件名
	safeFilename := sanitizeFilename(file.Filename)
	dstPath := filepath.Join(fileDir, safeFilename)

	// 4. 创建目标文件
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 5. 复制文件内容
	written, err := io.Copy(dst, src)
	if err != nil {
		// 上传失败，清理已创建的文件
		os.Remove(dstPath)
		return "", 0, fmt.Errorf("failed to save file: %w", err)
	}

	// 6. 返回相对路径（用于数据库存储）
	relativePath := filepath.Join(fmt.Sprintf("user_%d", userID), fileID, safeFilename)
	relativePath = filepath.ToSlash(relativePath) // 统一使用正斜杠

	logger.Info("File uploaded successfully",
		zap.String("path", relativePath),
		zap.Int64("size", written),
		zap.Uint("user_id", userID))

	return relativePath, written, nil
}

// Download 下载文件
func (s *LocalStorage) Download(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, filePath)

	// 安全检查：防止路径遍历攻击
	if !strings.HasPrefix(fullPath, s.basePath) {
		return nil, fmt.Errorf("invalid file path")
	}

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete 删除文件
func (s *LocalStorage) Delete(ctx context.Context, filePath string) error {
	fullPath := filepath.Join(s.basePath, filePath)

	// 安全检查
	if !strings.HasPrefix(fullPath, s.basePath) {
		return fmt.Errorf("invalid file path")
	}

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，视为删除成功
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	logger.Info("File deleted successfully", zap.String("path", filePath))
	return nil
}

// GetURL 获取文件访问URL
func (s *LocalStorage) GetURL(ctx context.Context, filePath string) (string, error) {
	// 本地存储使用静态文件服务
	return fmt.Sprintf("%s/%s", s.baseURL, filePath), nil
}

// GetSize 获取文件大小
func (s *LocalStorage) GetSize(ctx context.Context, filePath string) (int64, error) {
	fullPath := filepath.Join(s.basePath, filePath)

	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return info.Size(), nil
}

// Exists 检查文件是否存在
func (s *LocalStorage) Exists(ctx context.Context, filePath string) (bool, error) {
	fullPath := filepath.Join(s.basePath, filePath)

	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// sanitizeFilename 清理文件名，移除危险字符
func sanitizeFilename(filename string) string {
	// 替换路径分隔符
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "..", "_")

	// 移除其他危险字符
	dangerous := []string{"<", ">", ":", "\"", "|", "?", "*", "&", ";", "$"}
	for _, char := range dangerous {
		filename = strings.ReplaceAll(filename, char, "_")
	}

	// 如果文件名为空，使用时间戳
	if filename == "" {
		filename = fmt.Sprintf("file_%d", time.Now().UnixNano())
	}

	return filename
}
