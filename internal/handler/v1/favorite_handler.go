// internal/handler/v1/favorite_handler.go
package v1

import (
	"cloud-disk/internal/constant"
	"cloud-disk/internal/model/dto/request"
	"cloud-disk/internal/service/interfaces"
	"cloud-disk/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FavoriteHandler struct {
	favoriteService interfaces.FavoriteService
	validator       *validator.Validate
}

func NewFavoriteHandler(favoriteService interfaces.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
		validator:       validator.New(),
	}
}

// Add 添加收藏
// @Summary 添加收藏
// @Tags 收藏管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body request.FavoriteCreateRequest true "收藏信息"
// @Success 200 {object} response.Response{data=response.FavoriteResponse} "成功"
// @Router /api/v1/favorites [post]
func (h *FavoriteHandler) Add(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.FavoriteCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	favorite, err := h.favoriteService.Add(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Created(c, favorite)
}

// Remove 取消收藏
// @Summary 取消收藏
// @Tags 收藏管理
// @Security ApiKeyAuth
// @Produce json
// @Param file_id path string true "文件ID"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/favorites/{file_id} [delete]
func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("file_id")

	if err := h.favoriteService.Remove(c.Request.Context(), userID, fileID); err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetList 获取收藏列表
// @Summary 获取收藏列表
// @Tags 收藏管理
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=response.FavoriteListResponse} "成功"
// @Router /api/v1/favorites [get]
func (h *FavoriteHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.FavoriteListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	list, err := h.favoriteService.GetList(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, list)
}

// Check 检查是否已收藏
// @Summary 检查是否已收藏
// @Tags 收藏管理
// @Security ApiKeyAuth
// @Produce json
// @Param file_id query string true "文件ID"
// @Success 200 {object} response.Response{data=bool} "成功"
// @Router /api/v1/favorites/check [get]
func (h *FavoriteHandler) Check(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Query("file_id")

	if fileID == "" {
		response.BadRequest(c, constant.ErrParamInvalid)
		return
	}

	exists, err := h.favoriteService.IsFavorite(c.Request.Context(), userID, fileID)
	if err != nil {
		response.ErrorWithMessage(c, constant.ErrDatabase, err.Error())
		return
	}

	response.Success(c, exists)
}
