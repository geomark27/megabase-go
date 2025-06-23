package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string    `gorm:"size:100;not null" json:"name"`
	UserName      string    `gorm:"size:100;not null;uniqueIndex" json:"user_name"`
	Email         string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password      string    `gorm:"size:255;not null" json:"-"`
	RoleID        uint      `gorm:"not null" json:"role_id"`
	Role          Role      `gorm:"foreignKey:RoleID" json:"role"`
	RememberToken string    `gorm:"size:100;uniqueIndex" json:"-"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`
	LastLoginAt   time.Time `json:"last_login_at"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
