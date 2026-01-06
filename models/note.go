package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title" binding:"required"`
	Content   string         `gorm:"type:text" json:"content"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}
