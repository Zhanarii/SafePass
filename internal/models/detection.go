package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type DetectionStatus string

const (
	StatusCompliant DetectionStatus = "compliant"
	StatusViolation DetectionStatus = "violation"
	StatusWarning   DetectionStatus = "warning"
)

type PPEItem string

const (
	PPEHelmet  PPEItem = "helmet"
	PPEVest    PPEItem = "vest"
	PPEBoots   PPEItem = "boots"
	PPEGloves  PPEItem = "gloves"
	PPEGoggles PPEItem = "goggles"
	PPEMask    PPEItem = "mask"
)

type Detection struct {
	ID               uuid.UUID              `json:"id" db:"id"`
	CameraID         uuid.UUID              `json:"camera_id" db:"camera_id"`
	UserID           *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	AccessZoneID     *uuid.UUID             `json:"access_zone_id,omitempty" db:"access_zone_id"`
	Timestamp        time.Time              `json:"timestamp" db:"timestamp"`
	FrameURL         string                 `json:"frame_url" db:"frame_url"`
	DetectedPPE      pq.StringArray         `json:"detected_ppe" db:"detected_ppe"`
	MissingPPE       pq.StringArray         `json:"missing_ppe" db:"missing_ppe"`
	ConfidenceScores map[string]float64     `json:"confidence_scores" db:"confidence_scores"`
	BoundingBoxes    map[string]interface{} `json:"bounding_boxes" db:"bounding_boxes"`
	Status           DetectionStatus        `json:"status" db:"status"`
	FaceEmbedding    []byte                 `json:"-" db:"face_embedding"`
	ProcessingTimeMS int                    `json:"processing_time_ms" db:"processing_time_ms"`
	ModelVersion     string                 `json:"model_version" db:"model_version"`
}
