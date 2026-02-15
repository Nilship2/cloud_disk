// internal/model/entity/share.go
package entity

import (
	"time"

	"gorm.io/gorm"
)

// Share 分享实体
type Share struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Token         string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"token"` // 分享令牌
	FileID        string         `gorm:"type:varchar(36);not null;index" json:"file_id"`     // 文件ID
	UserID        uint           `gorm:"not null;index" json:"user_id"`                      // 分享用户ID
	Password      string         `gorm:"type:varchar(255)" json:"-"`                         // 访问密码（加密）
	ExpireTime    *time.Time     `json:"expire_time"`                                        // 过期时间
	MaxDownloads  int            `gorm:"default:0" json:"max_downloads"`                     // 最大下载次数（0表示无限制）
	DownloadCount int            `gorm:"default:0" json:"download_count"`                    // 已下载次数
	Status        int            `gorm:"default:1" json:"status"`                            // 状态：1-有效，2-已取消
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Share) TableName() string {
	return "shares"
}

// IsExpired 检查是否过期
func (s *Share) IsExpired() bool {
	if s.ExpireTime == nil {
		return false
	}
	return time.Now().After(*s.ExpireTime)
}

// IsDownloadLimitReached 检查是否达到下载限制
func (s *Share) IsDownloadLimitReached() bool {
	if s.MaxDownloads <= 0 {
		return false
	}
	return s.DownloadCount >= s.MaxDownloads
}

// CanAccess 检查是否可以访问
func (s *Share) CanAccess(password string) bool {
	// 检查状态
	if s.Status != 1 {
		return false
	}

	// 检查过期
	if s.IsExpired() {
		return false
	}

	// 检查下载限制
	if s.IsDownloadLimitReached() {
		return false
	}

	// 检查密码（如果有）
	if s.Password != "" {
		// 这里需要密码验证，将在service层实现
		return false
	}

	return true
}

// ShareResponse 分享响应DTO
type ShareResponse struct {
	ID            uint       `json:"id"`
	Token         string     `json:"token"`
	FileID        string     `json:"file_id"`
	Filename      string     `json:"filename"`       // 冗余字段，方便显示
	FileSize      int64      `json:"file_size"`      // 冗余字段
	FileSizeText  string     `json:"file_size_text"` // 格式化后的大小
	ExpireTime    *time.Time `json:"expire_time"`
	MaxDownloads  int        `json:"max_downloads"`
	DownloadCount int        `json:"download_count"`
	Status        int        `json:"status"`
	ShareLink     string     `json:"share_link"` // 完整分享链接
	CreatedAt     time.Time  `json:"created_at"`
}

// ShareDetailResponse 分享详情响应（包含文件信息）
type ShareDetailResponse struct {
	Token        string        `json:"token"`
	File         *FileResponse `json:"file"`
	ExpireTime   *time.Time    `json:"expire_time"`
	DownloadLeft int           `json:"download_left"` // 剩余下载次数
	NeedPassword bool          `json:"need_password"` // 是否需要密码
	CreatedAt    time.Time     `json:"created_at"`
}
