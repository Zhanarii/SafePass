package models

import (
	"time"

	"github.com/google/uuid"
)

type AccessLog struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	UserID         *uuid.UUID     `json:"user_id,omitempty" db:"user_id"`
	CameraID       uuid.UUID      `json:"camera_id" db:"camera_id"`
	AccessZoneID   *uuid.UUID     `json:"access_zone_id,omitempty" db:"access_zone_id"`
	DetectionID    *uuid.UUID     `json:"detection_id,omitempty" db:"detection_id"`
	Decision       AccessDecision `json:"decision" db:"decision"`
	Timestamp      time.Time      `json:"timestamp" db:"timestamp"`
	BadgeScanned   bool           `json:"badge_scanned" db:"badge_scanned"`
	OverrideBy     *uuid.UUID     `json:"override_by,omitempty" db:"override_by"`
	OverrideReason *string        `json:"override_reason,omitempty" db:"override_reason"`
}
