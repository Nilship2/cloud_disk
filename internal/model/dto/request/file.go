// internal/model/dto/request/file.go
package request

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	ParentID string `form:"parent_id" json:"parent_id" binding:"omitempty,uuid"`
}

// FileListRequest 文件列表请求
type FileListRequest struct {
	ParentID string `form:"parent_id" json:"parent_id" binding:"omitempty,uuid"`
	Page     int    `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
	OrderBy  string `form:"order_by" json:"order_by" binding:"omitempty,oneof=filename size created_at updated_at"`
	Order    string `form:"order" json:"order" binding:"omitempty,oneof=asc desc"`
}

// FileCreateFolderRequest 创建文件夹请求
type FileCreateFolderRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=255"`
	ParentID string `json:"parent_id" binding:"omitempty,uuid"`
}

// FileRenameRequest 重命名请求
type FileRenameRequest struct {
	NewName string `json:"new_name" binding:"required,min=1,max=255"`
}

// FileMoveRequest 移动文件请求
type FileMoveRequest struct {
	TargetParentID string `json:"target_parent_id" binding:"required,uuid"`
}

// FileBatchRequest 批量操作请求
type FileBatchRequest struct {
	FileIDs []string `json:"file_ids" binding:"required,min=1,dive,uuid"`
}

// FileInstantUploadRequest 秒传请求
type FileInstantUploadRequest struct {
	Hash     string `json:"hash" binding:"required"`
	Filename string `json:"filename" binding:"required"`
	ParentID string `json:"parent_id" binding:"omitempty,uuid"`
}
