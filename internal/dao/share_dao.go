// internal/dao/share_dao.go
package dao

import (
	"errors"
	"time"

	"cloud-disk/internal/model/entity"

	"gorm.io/gorm"
)

type ShareDAO struct {
	db *gorm.DB
}

func NewShareDAO(db *gorm.DB) *ShareDAO {
	return &ShareDAO{db: db}
}

// Create 创建分享
func (dao *ShareDAO) Create(share *entity.Share) error {
	return dao.db.Create(share).Error
}

// GetByToken 根据token获取分享
func (dao *ShareDAO) GetByToken(token string) (*entity.Share, error) {
	var share entity.Share
	err := dao.db.Where("token = ? AND deleted_at IS NULL", token).First(&share).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &share, err
}

// GetByID 根据ID获取分享
func (dao *ShareDAO) GetByID(id uint) (*entity.Share, error) {
	var share entity.Share
	err := dao.db.Where("id = ? AND deleted_at IS NULL", id).First(&share).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &share, err
}

// GetByFileID 根据文件ID获取分享
func (dao *ShareDAO) GetByFileID(fileID string) (*entity.Share, error) {
	var share entity.Share
	err := dao.db.Where("file_id = ? AND status = 1 AND deleted_at IS NULL", fileID).First(&share).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &share, err
}

// ListByUserID 获取用户分享列表
func (dao *ShareDAO) ListByUserID(userID uint, status int, page, pageSize int) ([]*entity.Share, int64, error) {
	var shares []*entity.Share
	var total int64

	db := dao.db.Model(&entity.Share{}).Where("user_id = ? AND deleted_at IS NULL", userID)

	if status > 0 {
		db = db.Where("status = ?", status)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&shares).Error; err != nil {
		return nil, 0, err
	}

	return shares, total, nil
}

// Update 更新分享
func (dao *ShareDAO) Update(share *entity.Share) error {
	return dao.db.Save(share).Error
}

// UpdateFields 更新指定字段
func (dao *ShareDAO) UpdateFields(id uint, fields map[string]interface{}) error {
	return dao.db.Model(&entity.Share{}).Where("id = ?", id).Updates(fields).Error
}

// IncrementDownloadCount 增加下载次数
func (dao *ShareDAO) IncrementDownloadCount(token string) error {
	return dao.db.Model(&entity.Share{}).Where("token = ?", token).
		Update("download_count", gorm.Expr("download_count + 1")).Error
}

// Delete 删除分享
func (dao *ShareDAO) Delete(id uint) error {
	return dao.db.Where("id = ?", id).Delete(&entity.Share{}).Error
}

// DeleteByFileID 删除文件的所有分享
func (dao *ShareDAO) DeleteByFileID(fileID string) error {
	return dao.db.Where("file_id = ?", fileID).Delete(&entity.Share{}).Error
}

// CleanExpired 清理过期分享
func (dao *ShareDAO) CleanExpired() error {
	return dao.db.Where("expire_time < ?", time.Now()).Delete(&entity.Share{}).Error
}
