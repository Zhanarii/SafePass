package models

import (
	"time"

	"github.com/google/uuid"
)

type Camera struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	LocationID    uuid.UUID  `json:"location_id" db:"location_id"`
	AccessZoneID  *uuid.UUID `json:"access_zone_id,omitempty" db:"access_zone_id"`
	Name          string     `json:"name" db:"name"`
	RTSPURL       string     `json:"rtsp_url" db:"rtsp_url"`
	Position      *string    `json:"position,omitempty" db:"position"`
	ViewingAngle  *int       `json:"viewing_angle,omitempty" db:"viewing_angle"`
	FPS           int        `json:"fps" db:"fps"`
	Resolution    string     `json:"resolution" db:"resolution"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	LastHeartbeat *time.Time `json:"last_heartbeat,omitempty" db:"last_heartbeat"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}
