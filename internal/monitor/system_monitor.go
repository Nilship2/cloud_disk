// internal/monitor/system_monitor.go
package monitor

import (
	"runtime"
	"time"

	"cloud-disk/internal/dao"
	"cloud-disk/internal/model/entity"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"gorm.io/gorm"
)

// SystemMonitor 系统监控
type SystemMonitor struct {
	db        *gorm.DB
	userDAO   *dao.UserDAO
	fileDAO   *dao.FileDAO
	shareDAO  *dao.ShareDAO
	startTime time.Time
}

// NewSystemMonitor 创建系统监控
func NewSystemMonitor(db *gorm.DB, userDAO *dao.UserDAO, fileDAO *dao.FileDAO, shareDAO *dao.ShareDAO) *SystemMonitor {
	return &SystemMonitor{
		db:        db,
		userDAO:   userDAO,
		fileDAO:   fileDAO,
		shareDAO:  shareDAO,
		startTime: time.Now(),
	}
}

// SystemStats 系统统计信息
type SystemStats struct {
	Uptime       string `json:"uptime"`        // 运行时间
	GoVersion    string `json:"go_version"`    // Go版本
	NumGoroutine int    `json:"num_goroutine"` // Goroutine数量

	// CPU信息
	CPUUsage float64 `json:"cpu_usage"` // CPU使用率

	// 内存信息
	MemoryTotal uint64  `json:"memory_total"` // 总内存
	MemoryUsed  uint64  `json:"memory_used"`  // 已用内存
	MemoryUsage float64 `json:"memory_usage"` // 内存使用率

	// 磁盘信息
	DiskTotal uint64  `json:"disk_total"` // 磁盘总空间
	DiskUsed  uint64  `json:"disk_used"`  // 磁盘已用空间
	DiskUsage float64 `json:"disk_usage"` // 磁盘使用率

	// 数据库统计
	UserCount    int64 `json:"user_count"`    // 用户总数
	FileCount    int64 `json:"file_count"`    // 文件总数
	ShareCount   int64 `json:"share_count"`   // 分享总数
	TotalStorage int64 `json:"total_storage"` // 总存储空间
	UsedStorage  int64 `json:"used_storage"`  // 已用存储空间
}

// GetSystemStats 获取系统统计信息
func (m *SystemMonitor) GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		Uptime:       time.Since(m.startTime).String(),
	}

	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(100*time.Millisecond, false)
	if err == nil && len(cpuPercent) > 0 {
		stats.CPUUsage = cpuPercent[0]
	}

	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		stats.MemoryTotal = memInfo.Total
		stats.MemoryUsed = memInfo.Used
		stats.MemoryUsage = memInfo.UsedPercent
	}

	// 获取磁盘信息
	diskInfo, err := disk.Usage(".")
	if err == nil {
		stats.DiskTotal = diskInfo.Total
		stats.DiskUsed = diskInfo.Used
		stats.DiskUsage = diskInfo.UsedPercent
	}

	// 获取数据库统计
	var userCount int64
	m.db.Model(&entity.User{}).Count(&userCount)
	stats.UserCount = userCount

	var fileCount int64
	m.db.Model(&entity.File{}).Count(&fileCount)
	stats.FileCount = fileCount

	var shareCount int64
	m.db.Model(&entity.Share{}).Count(&shareCount)
	stats.ShareCount = shareCount

	// 计算存储空间
	var totalStorage int64
	m.db.Model(&entity.User{}).Select("IFNULL(SUM(capacity), 0)").Scan(&totalStorage)
	stats.TotalStorage = totalStorage

	var usedStorage int64
	m.db.Model(&entity.User{}).Select("IFNULL(SUM(used), 0)").Scan(&usedStorage)
	stats.UsedStorage = usedStorage

	return stats, nil
}

// GetDBStats 获取数据库统计
func (m *SystemMonitor) GetDBStats() map[string]interface{} {
	sqlDB, err := m.db.DB()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// HealthCheck 健康检查
func (m *SystemMonitor) HealthCheck() map[string]interface{} {
	status := "healthy"

	// 检查数据库连接
	sqlDB, err := m.db.DB()
	if err != nil {
		status = "unhealthy"
	} else if err := sqlDB.Ping(); err != nil {
		status = "unhealthy"
	}

	return map[string]interface{}{
		"status":    status,
		"timestamp": time.Now(),
		"uptime":    time.Since(m.startTime).String(),
	}
}
