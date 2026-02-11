// internal/model/entity/user.go
package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Avatar    string         `gorm:"type:varchar(500)" json:"avatar"`
	Bio       string         `gorm:"type:varchar(500)" json:"bio"`
	Capacity  int64          `gorm:"not null;default:10737418240" json:"capacity"`
	Used      int64          `gorm:"not null;default:0" json:"used"`
	IsActive  bool           `gorm:"not null;default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// UserResponse 用户响应DTO
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Capacity  int64     `json:"capacity"`
	Used      int64     `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 转换为响应DTO
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Capacity:  u.Capacity,
		Used:      u.Used,
		CreatedAt: u.CreatedAt,
	}
}
