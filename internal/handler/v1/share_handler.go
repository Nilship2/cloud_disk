// internal/handler/v1/share_handler.go
package v1

import (
	"strconv"

	"cloud-disk/internal/constant"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ShareHandler struct {
	shareService interfaces.ShareService
	validator    *validator.Validate
}

func NewShareHandler(shareService interfaces.ShareService) *ShareHandler {
	return &ShareHandler{
		shareService: shareService,
		validator:    validator.New(),
	}
}

// Create 创建分享
// @Summary 创建分享
// @Description 创建文件分享链接
// @Tags 分享管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body request.ShareCreateRequest true "分享信息"
// @Success 201 {object} response.Response{data=response.ShareResponse} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "文件不存在"
// @Router /shares [post]
func (h *ShareHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.ShareCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	share, err := h.shareService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Created(c, share)
}

// Cancel 取消分享
// @Summary 取消分享
// @Tags 分享管理
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "分享ID"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/shares/{id} [delete]
func (h *ShareHandler) Cancel(c *gin.Context) {
	userID := c.GetUint("user_id")

	shareID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.shareService.Cancel(c.Request.Context(), userID, uint(shareID)); err != nil {
		response.ErrorWithMessage(c, constant.ErrShareNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

// Update 更新分享
// @Summary 更新分享
// @Tags 分享管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "分享ID"
// @Param request body request.ShareUpdateRequest true "更新信息"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/shares/{id} [put]
func (h *ShareHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")

	shareID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	var req request.ShareUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	if err := h.shareService.Update(c.Request.Context(), userID, uint(shareID), &req); err != nil {
		response.ErrorWithMessage(c, constant.ErrShareNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetDetail 获取分享详情
// @Summary 获取分享详情
// @Tags 分享管理
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "分享ID"
// @Success 200 {object} response.Response{data=response.ShareResponse} "成功"
// @Router /api/v1/shares/{id} [get]
func (h *ShareHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")

	shareID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	share, err := h.shareService.GetDetail(c.Request.Context(), userID, uint(shareID))
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrShareNotFound, err.Error())
		return
	}

	response.Success(c, share)
}

// GetList 获取分享列表
// @Summary 获取分享列表
// @Tags 分享管理
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query int false "状态"
// @Success 200 {object} response.Response{data=response.ShareListResponse} "成功"
// @Router /api/v1/shares [get]
func (h *ShareHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.ShareListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	list, err := h.shareService.GetList(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, list)
}

// Access 访问分享
// @Summary 访问分享
// @Description 通过token访问分享的文件信息
// @Tags 分享管理
// @Produce json
// @Param token path string true "分享令牌"
// @Param password query string false "访问密码"
// @Success 200 {object} response.Response{data=response.ShareDetailResponse} "成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 403 {object} response.ErrorResponse "密码错误"
// @Failure 404 {object} response.ErrorResponse "分享不存在"
// @Router /s/{token} [get]
func (h *ShareHandler) Access(c *gin.Context) {
	token := c.Param("token")
	password := c.Query("password")

	detail, err := h.shareService.AccessByToken(c.Request.Context(), token, password)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrShareNotFound, err.Error())
		return
	}

	response.Success(c, detail)
}

// Download 下载分享文件
// @Summary 下载分享文件
// @Tags 分享管理
// @Produce octet-stream
// @Param token path string true "分享令牌"
// @Param password query string false "访问密码"
// @Success 302 {string} string "重定向到文件"
// @Router /s/{token}/download [get]
func (h *ShareHandler) Download(c *gin.Context) {
	token := c.Param("token")
	password := c.Query("password")

	url, err := h.shareService.DownloadShare(c.Request.Context(), token, password)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrShareNotFound, err.Error())
		return
	}

	c.Redirect(302, url)
}
