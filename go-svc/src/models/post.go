package models

import "time"

type Post struct {
	ID            uint `gorm:"primaryKey"`
	Title         string
	Content       string
	IsActive      bool
	CreatedAt     time.Time
	CreatedBy     string
	LastUpdatedAt *time.Time
	LastUpdatedBy *string
}
