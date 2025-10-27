package models

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Address   string     `json:"address" db:"address"`
	Location  *string    `json:"location,omitempty" db:"location"`
	Type      string     `json:"type" db:"type"`
	CompanyID *uuid.UUID `json:"company_id,omitempty" db:"company_id"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}
