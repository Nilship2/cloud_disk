// internal/service/interfaces/file_service.go
package interfaces

import (
	"context"
	"mime/multipart"

	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
)

type FileService interface {
	// 上传文件
	Upload(ctx context.Context, userID uint, file *multipart.FileHeader, parentID string) (*response.FileResponse, error)

	// 秒传文件
	InstantUpload(ctx context.Context, userID uint, hash, filename string, parentID string) (*response.FileResponse, error)

	// 下载文件
	Download(ctx context.Context, userID uint, fileID string) (string, error)

	// 获取文件列表
	GetList(ctx context.Context, userID uint, req *request.FileListRequest) (*response.FileListResponse, error)

	// 获取文件详情
	GetDetail(ctx context.Context, userID uint, fileID string) (*response.FileResponse, error)

	// 删除文件（移入回收站）
	Delete(ctx context.Context, userID uint, fileID string) error

	// 批量删除
	BatchDelete(ctx context.Context, userID uint, fileIDs []string) error

	// 创建文件夹
	CreateFolder(ctx context.Context, userID uint, name, parentID string) (*response.FileResponse, error) // 改为 response.FileResponse

	// 重命名文件/文件夹
	Rename(ctx context.Context, userID uint, fileID, newName string) error

	// 移动文件/文件夹
	Move(ctx context.Context, userID uint, fileID, targetParentID string) error

	// 检查文件是否存在
	CheckExists(ctx context.Context, userID uint, hash string) (bool, *response.FileResponse) // 已修改
}
