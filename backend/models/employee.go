package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name"`
	Position   string         `json:"position"`
	BadgeID    string         `json:"badge_id" gorm:"uniqueIndex"`
	Status     string         `json:"status"` // active, inactive
	PhotoURL   string         `json:"photo_url"`
	Blocked    bool           `json:"blocked" gorm:"default:false"`
	LastSeenAt *time.Time     `json:"last_seen_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
