// internal/handler/v1/monitor_handler.go
package v1

import (
	"cloud-disk/internal/monitor"
	"cloud-disk/pkg/response"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	monitor *monitor.SystemMonitor
}

func NewMonitorHandler(monitor *monitor.SystemMonitor) *MonitorHandler {
	return &MonitorHandler{
		monitor: monitor,
	}
}

// GetSystemStats 获取系统统计信息
// @Summary 获取系统统计信息
// @Tags 监控
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=monitor.SystemStats} "成功"
// @Router /api/v1/monitor/stats [get]
func (h *MonitorHandler) GetSystemStats(c *gin.Context) {
	// 仅管理员可访问（这里简化，只检查用户存在）
	if c.GetUint("user_id") == 0 {
		response.Unauthorized(c)
		return
	}

	stats, err := h.monitor.GetSystemStats()
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, stats)
}

// GetDBStats 获取数据库统计
// @Summary 获取数据库统计
// @Tags 监控
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Router /api/v1/monitor/db [get]
func (h *MonitorHandler) GetDBStats(c *gin.Context) {
	if c.GetUint("user_id") == 0 {
		response.Unauthorized(c)
		return
	}

	stats := h.monitor.GetDBStats()
	response.Success(c, stats)
}

// Health 健康检查
// @Summary 健康检查
// @Tags 监控
// @Produce json
// @Success 200 {object} map[string]interface{} "健康状态"
// @Router /health [get]
func (h *MonitorHandler) Health(c *gin.Context) {
	health := h.monitor.HealthCheck()
	c.JSON(200, health)
}
