package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name          string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	DisplayName   string    `gorm:"size:100;not null" json:"display_name"`
	Description   string    `gorm:"type:text" json:"description"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Users         []User    `gorm:"foreignKey:RoleID" json:"-"`
}