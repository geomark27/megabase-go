package models

import (
	"gorm.io/gorm"
)

type Company struct{
	gorm.Model
	Name 		string `gorm:"size:100;not null" json:"name"`
	Host 		string `gorm:"size:100;not null" json:"host"`
	Database 	string `gorm:"size:100;not null" json:"database"`
	User 		string `gorm:"size:100;not null" json:"user"`
	Password	string `gorm:"size:100;not null" json:"password"`
	IsActive 	bool `gorm:"default:true" json:"is_active"`
}