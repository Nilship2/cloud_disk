// internal/service/interfaces/favorite_service.go
package interfaces

import (
	"context"

	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
)

type FavoriteService interface {
	// 添加收藏
	Add(ctx context.Context, userID uint, req *request.FavoriteCreateRequest) (*response.FavoriteResponse, error)

	// 取消收藏
	Remove(ctx context.Context, userID uint, fileID string) error

	// 获取收藏列表
	GetList(ctx context.Context, userID uint, req *request.FavoriteListRequest) (*response.FavoriteListResponse, error)

	// 检查是否已收藏
	IsFavorite(ctx context.Context, userID uint, fileID string) (bool, error)
}
