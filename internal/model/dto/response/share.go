// internal/model/dto/response/share.go
package response

import (
	"time"
)

// ShareResponse 分享响应
type ShareResponse struct {
	ID            uint       `json:"id"`
	Token         string     `json:"token"`
	FileID        string     `json:"file_id"`
	Filename      string     `json:"filename"`
	FileSize      int64      `json:"file_size"`
	FileSizeText  string     `json:"file_size_text"`
	ExpireTime    *time.Time `json:"expire_time"`
	MaxDownloads  int        `json:"max_downloads"`
	DownloadCount int        `json:"download_count"`
	Status        int        `json:"status"`
	ShareLink     string     `json:"share_link"`
	CreatedAt     time.Time  `json:"created_at"`
}

// ShareDetailResponse 分享详情响应
type ShareDetailResponse struct {
	Token        string        `json:"token"`
	File         *FileResponse `json:"file"`
	ExpireTime   *time.Time    `json:"expire_time"`
	DownloadLeft int           `json:"download_left"`
	NeedPassword bool          `json:"need_password"`
	CreatedAt    time.Time     `json:"created_at"`
}

// ShareListResponse 分享列表响应
type ShareListResponse struct {
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	Shares   []*ShareResponse `json:"shares"`
}
