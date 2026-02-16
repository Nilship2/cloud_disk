// internal/service/interfaces/trash_service.go
package interfaces

import (
	"context"

	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
)

type TrashService interface {
	// 获取回收站列表
	GetList(ctx context.Context, userID uint, req *request.TrashListRequest) (*response.TrashListResponse, error)

	// 恢复文件
	Restore(ctx context.Context, userID uint, fileID string) error

	// 批量恢复
	BatchRestore(ctx context.Context, userID uint, req *request.TrashBatchRequest) error

	// 彻底删除
	Delete(ctx context.Context, userID uint, fileID string) error

	// 批量彻底删除
	BatchDelete(ctx context.Context, userID uint, req *request.TrashBatchRequest) error

	// 清空回收站
	CleanAll(ctx context.Context, userID uint) (*response.TrashCleanResponse, error)

	// 获取回收站统计信息
	GetStats(ctx context.Context, userID uint) (int64, int64, error)
}
