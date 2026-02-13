// internal/handler/v1/file_handler.go
package v1

import (
	"cloud-disk/internal/constant"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FileHandler struct {
	fileService interfaces.FileService
	validator   *validator.Validate
}

func NewFileHandler(fileService interfaces.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		validator:   validator.New(),
	}
}

// Upload 上传文件
// @Summary 上传文件
// @Tags 文件管理
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param parent_id formData string false "父目录ID"
// @Success 200 {object} response.Response{data=entity.FileResponse} "成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Router /api/v1/files/upload [post]
func (h *FileHandler) Upload(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	// 获取父目录ID
	parentID := c.PostForm("parent_id")

	// 上传文件
	fileResp, err := h.fileService.Upload(c.Request.Context(), userID, file, parentID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrFileUploadFailed, err.Error())
		return
	}

	response.Created(c, fileResp)
}

// InstantUpload 秒传文件
// @Summary 秒传文件
// @Tags 文件管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body request.FileInstantUploadRequest true "秒传信息"
// @Success 200 {object} response.Response{data=entity.FileResponse} "成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Router /api/v1/files/instant [post]
func (h *FileHandler) InstantUpload(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.FileInstantUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	fileResp, err := h.fileService.InstantUpload(c.Request.Context(), userID, req.Hash, req.Filename, req.ParentID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrFileNotFound, err.Error())
		return
	}

	response.Success(c, fileResp)
}

// Download 下载文件
// @Summary 下载文件
// @Tags 文件管理
// @Security ApiKeyAuth
// @Produce octet-stream
// @Param id path string true "文件ID"
// @Success 200 {file} binary "文件"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Router /api/v1/files/{id}/download [get]
func (h *FileHandler) Download(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	url, err := h.fileService.Download(c.Request.Context(), userID, fileID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrFileNotFound, err.Error())
		return
	}

	c.Redirect(302, url)
}

// GetList 获取文件列表
// @Summary 获取文件列表
// @Tags 文件管理
// @Security ApiKeyAuth
// @Produce json
// @Param parent_id query string false "父目录ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param order_by query string false "排序字段"
// @Param order query string false "排序方式"
// @Success 200 {object} response.Response{data=response.FileListResponse} "成功"
// @Router /api/v1/files [get]
func (h *FileHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.FileListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	list, err := h.fileService.GetList(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, list)
}

// GetDetail 获取文件详情
// @Summary 获取文件详情
// @Tags 文件管理
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "文件ID"
// @Success 200 {object} response.Response{data=entity.FileResponse} "成功"
// @Router /api/v1/files/{id} [get]
func (h *FileHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	file, err := h.fileService.GetDetail(c.Request.Context(), userID, fileID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrFileNotFound, err.Error())
		return
	}

	response.Success(c, file)
}

// Delete 删除文件（移入回收站）
// @Summary 删除文件
// @Tags 文件管理
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "文件ID"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/files/{id} [delete]
func (h *FileHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	if err := h.fileService.Delete(c.Request.Context(), userID, fileID); err != nil {
		response.ErrorWithMessage(c, constant.ErrFileNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchDelete 批量删除文件
// @Summary 批量删除文件
// @Tags 文件管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body request.FileBatchRequest true "文件ID列表"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/files/batch [delete]
func (h *FileHandler) BatchDelete(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.FileBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.fileService.BatchDelete(c.Request.Context(), userID, req.FileIDs); err != nil {
		response.ErrorWithMessage(c, constant.ErrFileNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateFolder 创建文件夹
// @Summary 创建文件夹
// @Tags 文件管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body request.FileCreateFolderRequest true "文件夹信息"
// @Success 200 {object} response.Response{data=entity.FileResponse} "成功"
// @Router /api/v1/folders [post]
func (h *FileHandler) CreateFolder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.FileCreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	folder, err := h.fileService.CreateFolder(c.Request.Context(), userID, req.Name, req.ParentID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Created(c, folder)
}

// Rename 重命名文件/文件夹
// @Summary 重命名文件
// @Tags 文件管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "文件ID"
// @Param request body request.FileRenameRequest true "新名称"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/files/{id}/rename [put]
func (h *FileHandler) Rename(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	var req request.FileRenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.fileService.Rename(c.Request.Context(), userID, fileID, req.NewName); err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, nil)
}

// Move 移动文件/文件夹
// @Summary 移动文件
// @Tags 文件管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "文件ID"
// @Param request body request.FileMoveRequest true "目标目录"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/files/{id}/move [put]
func (h *FileHandler) Move(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	var req request.FileMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.fileService.Move(c.Request.Context(), userID, fileID, req.TargetParentID); err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, nil)
}

// CheckExists 检查文件是否存在（秒传）
// @Summary 检查文件是否存在
// @Tags 文件管理
// @Security ApiKeyAuth
// @Produce json
// @Param hash query string true "文件哈希"
// @Success 200 {object} response.Response{data=entity.FileResponse} "成功"
// @Router /api/v1/files/check [get]
func (h *FileHandler) CheckExists(c *gin.Context) {
	userID := c.GetUint("user_id")
	hash := c.Query("hash")

	if hash == "" {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	exists, file := h.fileService.CheckExists(c.Request.Context(), userID, hash)
	if exists {
		response.Success(c, file)
		return
	}

	response.Success(c, nil)
}
