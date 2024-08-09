package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Files    []File `gorm:"foreignKey:UserID"`
}

type File struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Size      int64
	Type      string
	UserID    uint
	UploadURL string
}