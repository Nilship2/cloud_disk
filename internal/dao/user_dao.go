// internal/dao/user_dao.go
package dao

import (
	"errors"
	"time"

	"cloud-disk/internal/model/entity"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

// Create 创建用户
func (dao *UserDAO) Create(user *entity.User) error {
	return dao.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (dao *UserDAO) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	err := dao.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByUsername 根据用户名获取用户
func (dao *UserDAO) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := dao.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByEmail 根据邮箱获取用户
func (dao *UserDAO) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := dao.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByUsernameOrEmail 根据用户名或邮箱获取用户
func (dao *UserDAO) GetByUsernameOrEmail(account string) (*entity.User, error) {
	var user entity.User
	err := dao.db.Where("(username = ? OR email = ?) AND deleted_at IS NULL", account, account).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Update 更新用户信息
func (dao *UserDAO) Update(user *entity.User) error {
	return dao.db.Save(user).Error
}

// UpdateFields 更新指定字段
func (dao *UserDAO) UpdateFields(id uint, fields map[string]interface{}) error {
	return dao.db.Model(&entity.User{}).Where("id = ?", id).Updates(fields).Error
}

// UpdateLastLogin 更新最后登录时间
func (dao *UserDAO) UpdateLastLogin(id uint) error {
	return dao.db.Model(&entity.User{}).Where("id = ?", id).
		Update("last_login", time.Now()).Error
}

// UpdatePassword 更新密码
func (dao *UserDAO) UpdatePassword(id uint, password string) error {
	return dao.db.Model(&entity.User{}).Where("id = ?", id).
		Update("password", password).Error
}

// UpdateStorageUsed 更新已用空间
func (dao *UserDAO) UpdateStorageUsed(id uint, size int64) error {
	return dao.db.Model(&entity.User{}).Where("id = ?", id).
		Update("used", gorm.Expr("used + ?", size)).Error
}

// CheckStorage 检查存储空间是否足够
func (dao *UserDAO) CheckStorage(id uint, size int64) (bool, error) {
	var user entity.User
	err := dao.db.Select("capacity", "used").Where("id = ?", id).First(&user).Error
	if err != nil {
		return false, err
	}
	return user.Used+size <= user.Capacity, nil
}

// Delete 软删除用户
func (dao *UserDAO) Delete(id uint) error {
	return dao.db.Where("id = ?", id).Delete(&entity.User{}).Error
}

// Exists 检查用户名或邮箱是否存在
func (dao *UserDAO) Exists(username, email string) (bool, error) {
	var count int64
	err := dao.db.Model(&entity.User{}).
		Where("(username = ? OR email = ?) AND deleted_at IS NULL", username, email).
		Count(&count).Error
	return count > 0, err
}
