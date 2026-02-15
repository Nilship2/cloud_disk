// internal/model/dto/request/share.go
package request

// ShareCreateRequest 创建分享请求
type ShareCreateRequest struct {
	FileID       string `json:"file_id" binding:"required,uuid"`
	Password     string `json:"password" binding:"omitempty,min=4,max=20"`        // 访问密码
	ExpireDays   int    `json:"expire_days" binding:"omitempty,min=1,max=30"`     // 过期天数
	MaxDownloads int    `json:"max_downloads" binding:"omitempty,min=1,max=1000"` // 最大下载次数
}

// ShareUpdateRequest 更新分享请求
type ShareUpdateRequest struct {
	Password     *string `json:"password" binding:"omitempty,min=4,max=20"`
	ExpireDays   *int    `json:"expire_days" binding:"omitempty,min=1,max=30"`
	MaxDownloads *int    `json:"max_downloads" binding:"omitempty,min=1,max=1000"`
	Status       *int    `json:"status" binding:"omitempty,oneof=1 2"`
}

// ShareAccessRequest 访问分享请求
type ShareAccessRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"omitempty"`
}

// ShareListRequest 分享列表请求
type ShareListRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=50"`
	Status   int `form:"status" json:"status" binding:"omitempty,oneof=0 1 2"` // 0-全部，1-有效，2-已取消
}
