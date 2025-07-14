package models

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    *uint          `json:"user_id" gorm:"index"`
	Author    *string        `json:"author,omitempty" gorm:"type:varchar(100);"`
	Title     string         `json:"title" gorm:"type:varchar(255);not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Comments  []Comment      `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
