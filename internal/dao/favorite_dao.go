// internal/dao/favorite_dao.go
package dao

import (
	"errors"

	"cloud-disk/internal/model/entity"

	"gorm.io/gorm"
)

type FavoriteDAO struct {
	db *gorm.DB
}

func NewFavoriteDAO(db *gorm.DB) *FavoriteDAO {
	return &FavoriteDAO{db: db}
}

// Create 创建收藏
func (dao *FavoriteDAO) Create(favorite *entity.Favorite) error {
	return dao.db.Create(favorite).Error
}

// Delete 取消收藏
func (dao *FavoriteDAO) Delete(userID uint, fileID string) error {
	return dao.db.Where("user_id = ? AND file_id = ?", userID, fileID).Delete(&entity.Favorite{}).Error
}

// GetByUserAndFile 获取用户的某个收藏
func (dao *FavoriteDAO) GetByUserAndFile(userID uint, fileID string) (*entity.Favorite, error) {
	var favorite entity.Favorite
	err := dao.db.Where("user_id = ? AND file_id = ?", userID, fileID).First(&favorite).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &favorite, err
}

// ListByUserID 获取用户收藏列表
func (dao *FavoriteDAO) ListByUserID(userID uint, page, pageSize int) ([]*entity.Favorite, int64, error) {
	var favorites []*entity.Favorite
	var total int64

	db := dao.db.Model(&entity.Favorite{}).Where("user_id = ?", userID)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页，预加载文件信息
	offset := (page - 1) * pageSize
	if err := db.Preload("File").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&favorites).Error; err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// CountByUserID 统计用户收藏数量
func (dao *FavoriteDAO) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := dao.db.Model(&entity.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// IsFavorite 检查是否已收藏
func (dao *FavoriteDAO) IsFavorite(userID uint, fileID string) (bool, error) {
	var count int64
	err := dao.db.Model(&entity.Favorite{}).Where("user_id = ? AND file_id = ?", userID, fileID).Count(&count).Error
	return count > 0, err
}
