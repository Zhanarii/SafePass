package models

import (
	"time"

	"github.com/google/uuid"
)

type ViolationSeverity string

const (
	SeverityLow      ViolationSeverity = "low"
	SeverityMedium   ViolationSeverity = "medium"
	SeverityHigh     ViolationSeverity = "high"
	SeverityCritical ViolationSeverity = "critical"
)

type AccessDecision string

const (
	AccessAllowed      AccessDecision = "allowed"
	AccessDenied       AccessDecision = "denied"
	AccessManualReview AccessDecision = "manual_review"
)

type Violation struct {
	ID             uuid.UUID         `json:"id" db:"id"`
	DetectionID    uuid.UUID         `json:"detection_id" db:"detection_id"`
	UserID         *uuid.UUID        `json:"user_id,omitempty" db:"user_id"`
	CameraID       uuid.UUID         `json:"camera_id" db:"camera_id"`
	AccessZoneID   *uuid.UUID        `json:"access_zone_id,omitempty" db:"access_zone_id"`
	ViolationType  string            `json:"violation_type" db:"violation_type"`
	Severity       ViolationSeverity `json:"severity" db:"severity"`
	Description    string            `json:"description" db:"description"`
	SnapshotURL    string            `json:"snapshot_url" db:"snapshot_url"`
	VideoURL       string            `json:"video_url,omitempty" db:"video_url"`
	AccessDecision AccessDecision    `json:"access_decision" db:"access_decision"`
	NotifiedAt     *time.Time        `json:"notified_at,omitempty" db:"notified_at"`
	AcknowledgedBy *uuid.UUID        `json:"acknowledged_by,omitempty" db:"acknowledged_by"`
	AcknowledgedAt *time.Time        `json:"acknowledged_at,omitempty" db:"acknowledged_at"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
}
