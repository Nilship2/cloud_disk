// internal/dao/file_dao.go
package dao

import (
	"errors"

	"cloud-disk/internal/model/entity"

	"gorm.io/gorm"
)

type FileDAO struct {
	db *gorm.DB
}

func NewFileDAO(db *gorm.DB) *FileDAO {
	return &FileDAO{db: db}
}

// Create 创建文件记录
func (dao *FileDAO) Create(file *entity.File) error {
	return dao.db.Create(file).Error
}

// GetByID 根据ID获取文件
func (dao *FileDAO) GetByID(id string) (*entity.File, error) {
	var file entity.File
	err := dao.db.Where("id = ? AND deleted_at IS NULL", id).First(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &file, err
}

// GetByUserAndPath 根据用户ID和路径获取文件
func (dao *FileDAO) GetByUserAndPath(userID uint, path string) (*entity.File, error) {
	var file entity.File
	err := dao.db.Where("user_id = ? AND path = ? AND deleted_at IS NULL", userID, path).First(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &file, err
}

// GetByHash 根据哈希值获取文件（用于秒传）
func (dao *FileDAO) GetByHash(hash string) (*entity.File, error) {
	var file entity.File
	err := dao.db.Where("hash = ? AND deleted_at IS NULL", hash).First(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &file, err
}

// ListByUserID 获取用户文件列表
func (dao *FileDAO) ListByUserID(userID uint, parentID *string, page, pageSize int, orderBy, order string) ([]*entity.File, int64, error) {
	var files []*entity.File
	var total int64

	db := dao.db.Model(&entity.File{}).Where("user_id = ? AND deleted_at IS NULL", userID)

	if parentID == nil {
		db = db.Where("parent_id IS NULL")
	} else {
		db = db.Where("parent_id = ?", parentID)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderStr := orderBy + " " + order
	// 分页
	offset := (page - 1) * pageSize

	if err := db.Order(orderStr).Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// Update 更新文件信息
func (dao *FileDAO) Update(file *entity.File) error {
	return dao.db.Save(file).Error
}

// UpdateFields 更新指定字段
func (dao *FileDAO) UpdateFields(id string, fields map[string]interface{}) error {
	return dao.db.Model(&entity.File{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 软删除文件（移入回收站）
func (dao *FileDAO) Delete(id string) error {
	return dao.db.Where("id = ?", id).Delete(&entity.File{}).Error
}

// BatchDelete 批量软删除
func (dao *FileDAO) BatchDelete(ids []string) error {
	return dao.db.Where("id IN (?)", ids).Delete(&entity.File{}).Error
}

// Restore 从回收站恢复
func (dao *FileDAO) Restore(id string) error {
	return dao.db.Model(&entity.File{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

// BatchRestore 批量恢复
func (dao *FileDAO) BatchRestore(ids []string) error {
	return dao.db.Model(&entity.File{}).Where("id IN (?)", ids).Update("deleted_at", nil).Error
}

// GetTrashList 获取回收站列表
func (dao *FileDAO) GetTrashList(userID uint, page, pageSize int) ([]*entity.File, int64, error) {
	var files []*entity.File
	var total int64

	db := dao.db.Model(&entity.File{}).Where("user_id = ? AND deleted_at IS NOT NULL", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("deleted_at DESC").Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// CountByUserID 统计用户文件数量
func (dao *FileDAO) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := dao.db.Model(&entity.File{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&count).Error
	return count, err
}

// SumSizeByUserID 统计用户文件总大小
func (dao *FileDAO) SumSizeByUserID(userID uint) (int64, error) {
	var sum int64
	err := dao.db.Model(&entity.File{}).Where("user_id = ? AND deleted_at IS NULL", userID).Select("IFNULL(SUM(size), 0)").Scan(&sum).Error
	return sum, err
}
