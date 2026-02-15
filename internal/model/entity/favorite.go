// internal/model/entity/favorite.go
package entity

import (
	"time"
)

// Favorite 收藏实体
type Favorite struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_file,unique" json:"user_id"`                  // 用户ID
	FileID    string    `gorm:"type:varchar(36);not null;index:idx_user_file,unique" json:"file_id"` // 文件ID
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

// FavoriteResponse 收藏响应DTO
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
