package models

import (
	"time"

	"gorm.io/gorm"
)

type Note struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title" binding:"required"`
	Content   string         `gorm:"type:text" json:"content"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Name      string         `gorm:"size:255" json:"name"`
	Notes     []Note         `gorm:"foreignKey:UserID" json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
