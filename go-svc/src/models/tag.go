package models

import "time"

type Tag struct {
	ID            uint `gorm:"primaryKey"`
	Label         string
	IsActive      bool
	CreatedAt     time.Time
	CreatedBy     string
	LastUpdatedAt *time.Time
	LastUpdatedBy *string
}
