package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AccessZone struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	LocationID  uuid.UUID              `json:"location_id" db:"location_id"`
	CameraID    *uuid.UUID             `json:"camera_id,omitempty" db:"camera_id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	RequiredPPE pq.StringArray         `json:"required_ppe" db:"required_ppe"`
	DangerLevel ViolationSeverity      `json:"danger_level" db:"danger_level"`
	AccessRules map[string]interface{} `json:"access_rules,omitempty" db:"access_rules"`
	IsActive    bool                   `json:"is_active" db:"is_active"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}
