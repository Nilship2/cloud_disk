// internal/model/dto/response/favorite.go
package response

import (
	"time"
)

// FavoriteResponse 收藏响应
type FavoriteResponse struct {
	ID        uint          `json:"id"`
	File      *FileResponse `json:"file"`
	CreatedAt time.Time     `json:"created_at"`
}

// FavoriteListResponse 收藏列表响应
type FavoriteListResponse struct {
	Total     int64               `json:"total"`
	Page      int                 `json:"page"`
	PageSize  int                 `json:"page_size"`
	Favorites []*FavoriteResponse `json:"favorites"`
}
