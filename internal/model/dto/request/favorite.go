// internal/model/dto/request/favorite.go
package request

// FavoriteCreateRequest 创建收藏请求
type FavoriteCreateRequest struct {
	FileID string `json:"file_id" binding:"required,uuid"`
}

// FavoriteListRequest 收藏列表请求
type FavoriteListRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=50"`
}
