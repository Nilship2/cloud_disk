// internal/service/interfaces/share_service.go
package interfaces

import (
	"context"

	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/dto/response"
)

type ShareService interface {
	// 创建分享
	Create(ctx context.Context, userID uint, req *request.ShareCreateRequest) (*response.ShareResponse, error)

	// 取消分享
	Cancel(ctx context.Context, userID uint, shareID uint) error

	// 更新分享
	Update(ctx context.Context, userID uint, shareID uint, req *request.ShareUpdateRequest) error

	// 获取分享详情
	GetDetail(ctx context.Context, userID uint, shareID uint) (*response.ShareResponse, error)

	// 获取用户分享列表
	GetList(ctx context.Context, userID uint, req *request.ShareListRequest) (*response.ShareListResponse, error)

	// 通过token访问分享
	AccessByToken(ctx context.Context, token, password string) (*response.ShareDetailResponse, error)

	// 下载分享文件
	DownloadShare(ctx context.Context, token, password string) (string, error)

	// 验证分享访问权限
	ValidateAccess(ctx context.Context, token, password string) (string, error)
}
