package models

import "time"

type PostsTag struct {
	ID            uint `gorm:"primaryKey"`
	PostId        uint
	TagId         uint
	IsActive      bool
	CreatedAt     time.Time
	CreatedBy     string
	LastUpdatedAt *time.Time
	LastUpdatedBy *string
}
