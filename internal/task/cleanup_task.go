// internal/task/cleanup_task.go
package task

import (
	"time"

	"cloud-disk/internal/dao"
	"cloud-disk/internal/model/entity"
	"cloud-disk/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CleanupTask 清理任务
type CleanupTask struct {
	shareDAO *dao.ShareDAO
	fileDAO  *dao.FileDAO
	db       *gorm.DB
	interval time.Duration
	stopChan chan struct{}
}

// NewCleanupTask 创建清理任务
func NewCleanupTask(db *gorm.DB, shareDAO *dao.ShareDAO, fileDAO *dao.FileDAO) *CleanupTask {
	return &CleanupTask{
		shareDAO: shareDAO,
		fileDAO:  fileDAO,
		db:       db,
		interval: 24 * time.Hour, // 默认每天执行一次
		stopChan: make(chan struct{}),
	}
}

// Start 启动清理任务
func (t *CleanupTask) Start() {
	logger.Info("Starting cleanup task", zap.Duration("interval", t.interval))

	go func() {
		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()

		// 立即执行一次
		t.run()

		for {
			select {
			case <-ticker.C:
				t.run()
			case <-t.stopChan:
				logger.Info("Cleanup task stopped")
				return
			}
		}
	}()
}

// Stop 停止清理任务
func (t *CleanupTask) Stop() {
	close(t.stopChan)
}

// run 执行清理
func (t *CleanupTask) run() {
	logger.Info("Running cleanup task...")

	// 1. 清理过期分享
	t.cleanExpiredShares()

	// 2. 清理过期回收站文件
	t.cleanExpiredTrash()

	logger.Info("Cleanup task completed")
}

// cleanExpiredShares 清理过期分享
func (t *CleanupTask) cleanExpiredShares() {
	// 软删除过期分享
	result := t.db.Where("expire_time < ?", time.Now()).Delete(&entity.Share{})
	if result.Error != nil {
		logger.Error("Failed to clean expired shares", zap.Error(result.Error))
	} else {
		logger.Info("Cleaned expired shares", zap.Int64("count", result.RowsAffected))
	}
}

// cleanExpiredTrash 清理过期回收站文件
func (t *CleanupTask) cleanExpiredTrash() {
	// 计算30天前的时间
	expireTime := time.Now().AddDate(0, 0, -30)

	// 获取要删除的文件列表（用于日志）
	var files []*entity.File
	if err := t.db.Unscoped().Where("deleted_at < ?", expireTime).Find(&files).Error; err != nil {
		logger.Error("Failed to query expired trash", zap.Error(err))
		return
	}

	if len(files) == 0 {
		logger.Info("No expired trash files found")
		return
	}

	// 计算总大小
	var totalSize int64
	for _, file := range files {
		totalSize += file.Size
	}

	// 彻底删除
	result := t.db.Unscoped().Where("deleted_at < ?", expireTime).Delete(&entity.File{})
	if result.Error != nil {
		logger.Error("Failed to clean expired trash", zap.Error(result.Error))
	} else {
		logger.Info("Cleaned expired trash",
			zap.Int64("count", result.RowsAffected),
			zap.Int64("size", totalSize),
			zap.String("size_text", entity.FormatFileSize(totalSize)))
	}
}

// SetInterval 设置执行间隔
func (t *CleanupTask) SetInterval(interval time.Duration) {
	t.interval = interval
}
