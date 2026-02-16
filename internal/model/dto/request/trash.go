// internal/model/dto/request/trash.go
package request

// TrashListRequest 回收站列表请求
type TrashListRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
}

// TrashBatchRequest 批量操作请求
type TrashBatchRequest struct {
	FileIDs []string `json:"file_ids" binding:"required,min=1,dive,uuid"`
}
