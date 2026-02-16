// internal/model/dto/response/trash.go
package response

import (
	"time"
)

// TrashItemResponse 回收站项响应
type TrashItemResponse struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	SizeText  string    `json:"size_text"`
	MimeType  string    `json:"mime_type"`
	Extension string    `json:"extension"`
	IsDir     bool      `json:"is_dir"`
	DeletedAt time.Time `json:"deleted_at"`
	ExpireIn  int       `json:"expire_in"` // 剩余过期天数（默认30天后永久删除）
}

// TrashListResponse 回收站列表响应
type TrashListResponse struct {
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
	Items    []*TrashItemResponse `json:"items"`
}

// TrashCleanResponse 清理响应
type TrashCleanResponse struct {
	CleanedCount   int64  `json:"cleaned_count"` // 清理的文件数
	FreedSpace     int64  `json:"freed_space"`   // 释放的空间
	FreedSpaceText string `json:"freed_space_text"`
}
