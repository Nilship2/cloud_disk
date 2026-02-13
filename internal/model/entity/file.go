// internal/model/entity/file.go
package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// File 文件实体
type File struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`      // UUID
	UserID    uint           `gorm:"not null;index" json:"user_id"`              // 所属用户
	Filename  string         `gorm:"type:varchar(255);not null" json:"filename"` // 文件名
	Path      string         `gorm:"type:varchar(500);not null" json:"path"`     // 存储路径
	Size      int64          `gorm:"not null;default:0" json:"size"`             // 文件大小(字节)
	Hash      string         `gorm:"type:varchar(64);index" json:"hash"`         // 文件哈希（用于秒传）
	MimeType  string         `gorm:"type:varchar(100)" json:"mime_type"`         // MIME类型
	Extension string         `gorm:"type:varchar(20)" json:"extension"`          // 文件扩展名
	IsDir     bool           `gorm:"not null;default:false" json:"is_dir"`       // 是否为目录
	ParentID  *string        `gorm:"type:varchar(36);index" json:"parent_id"`    // 父目录ID
	Status    int            `gorm:"not null;default:1" json:"status"`           // 状态：1-正常，2-冻结
	CreatedAt time.Time      `json:"created_at"`                                 // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                                 // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                             // 软删除时间（回收站）
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// FileResponse 文件响应DTO
type FileResponse struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	SizeText  string    `json:"size_text"` // 格式化后的文件大小
	MimeType  string    `json:"mime_type"`
	Extension string    `json:"extension"`
	IsDir     bool      `json:"is_dir"`
	ParentID  *string   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url,omitempty"` // 访问URL
}

// ToResponse 转换为响应DTO
func (f *File) ToResponse() *FileResponse {
	return &FileResponse{
		ID:        f.ID,
		Filename:  f.Filename,
		Path:      f.Path,
		Size:      f.Size,
		SizeText:  FormatFileSize(f.Size),
		MimeType:  f.MimeType,
		Extension: f.Extension,
		IsDir:     f.IsDir,
		ParentID:  f.ParentID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

// FormatFileSize 格式化文件大小
func FormatFileSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	ParentID string `form:"parent_id" binding:"omitempty,uuid"` // 父目录ID
}

// FileListRequest 文件列表请求
type FileListRequest struct {
	ParentID string `form:"parent_id" binding:"omitempty,uuid"`           // 父目录ID
	Page     int    `form:"page,default=1" binding:"min=1"`               // 页码
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"` // 每页数量
	OrderBy  string `form:"order_by,default=created_at"`                  // 排序字段
	Order    string `form:"order,default=desc"`                           // 排序方式
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Files    []*FileResponse `json:"files"`
}
