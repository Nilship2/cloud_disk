// internal/model/dto/response/file.go
package response

import (
	"time"
)

// FileResponse 文件响应
type FileResponse struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	SizeText  string    `json:"size_text"`
	MimeType  string    `json:"mime_type"`
	Extension string    `json:"extension"`
	IsDir     bool      `json:"is_dir"`
	ParentID  *string   `json:"parent_id"`
	URL       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Files    []*FileResponse `json:"files"`
}

// StorageInfoResponse 存储空间信息响应
type StorageInfoResponse struct {
	Capacity    int64   `json:"capacity"`
	Used        int64   `json:"used"`
	Available   int64   `json:"available"`
	UsageRate   float64 `json:"usage_rate"`
	FileCount   int64   `json:"file_count"`
	FolderCount int64   `json:"folder_count"`
}
