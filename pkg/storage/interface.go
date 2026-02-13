// pkg/storage/interface.go
package storage

import (
	"context"
	"io"
	"mime/multipart"
)

// Storage 存储接口
type Storage interface {
	// 上传文件
	Upload(ctx context.Context, file *multipart.FileHeader, userID uint, fileID string) (string, int64, error)

	// 下载文件
	Download(ctx context.Context, filePath string) (io.ReadCloser, error)

	// 删除文件
	Delete(ctx context.Context, filePath string) error

	// 获取文件访问URL
	GetURL(ctx context.Context, filePath string) (string, error)

	// 获取文件大小
	GetSize(ctx context.Context, filePath string) (int64, error)

	// 检查文件是否存在
	Exists(ctx context.Context, filePath string) (bool, error)
}

// FileInfo 文件信息
type FileInfo struct {
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	Name      string `json:"name"`
	MimeType  string `json:"mime_type"`
	Extension string `json:"extension"`
}
