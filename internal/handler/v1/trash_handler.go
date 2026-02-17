// internal/handler/v1/trash_handler.go
package v1

import (
	"cloud-disk/internal/constant"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/model/entity"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TrashHandler struct {
	trashService interfaces.TrashService
	validator    *validator.Validate
}

func NewTrashHandler(trashService interfaces.TrashService) *TrashHandler {
	return &TrashHandler{
		trashService: trashService,
		validator:    validator.New(),
	}
}

// GetList 获取回收站列表
// @Summary 获取回收站列表
// @Description 获取已删除的文件列表
// @Tags 回收站
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20) maximum(100)
// @Success 200 {object} response.Response{data=response.TrashListResponse} "成功"
// @Router /trash [get]
func (h *TrashHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.TrashListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	list, err := h.trashService.GetList(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, list)
}

// Restore 恢复文件
// @Summary 恢复文件
// @Description 从回收站恢复文件
// @Tags 回收站
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "文件ID"
// @Success 200 {object} response.Response "恢复成功"
// @Failure 404 {object} response.ErrorResponse "文件不存在"
// @Router /trash/{id}/restore [post]
func (h *TrashHandler) Restore(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	if err := h.trashService.Restore(c.Request.Context(), userID, fileID); err != nil {
		response.ErrorWithMessage(c, constant.ErrFileNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchRestore 批量恢复
// @Summary 批量恢复
// @Tags 回收站
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body request.TrashBatchRequest true "文件ID列表"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/trash/batch/restore [post]
func (h *TrashHandler) BatchRestore(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.TrashBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.trashService.BatchRestore(c.Request.Context(), userID, &req); err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 彻底删除
// @Summary 彻底删除
// @Tags 回收站
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "文件ID"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/trash/{id} [delete]
func (h *TrashHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	if err := h.trashService.Delete(c.Request.Context(), userID, fileID); err != nil {
		response.ErrorWithMessage(c, constant.ErrFileNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchDelete 批量彻底删除
// @Summary 批量彻底删除
// @Tags 回收站
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body request.TrashBatchRequest true "文件ID列表"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/trash/batch [delete]
func (h *TrashHandler) BatchDelete(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.TrashBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.trashService.BatchDelete(c.Request.Context(), userID, &req); err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, nil)
}

// CleanAll 清空回收站
// @Summary 清空回收站
// @Tags 回收站
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=response.TrashCleanResponse} "成功"
// @Router /api/v1/trash/clean [post]
func (h *TrashHandler) CleanAll(c *gin.Context) {
	userID := c.GetUint("user_id")

	resp, err := h.trashService.CleanAll(c.Request.Context(), userID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetStats 获取回收站统计
// @Summary 获取回收站统计
// @Tags 回收站
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Router /api/v1/trash/stats [get]
func (h *TrashHandler) GetStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, size, err := h.trashService.GetStats(c.Request.Context(), userID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, gin.H{
		"count":     count,
		"size":      size,
		"size_text": entity.FormatFileSize(size),
	})
}
