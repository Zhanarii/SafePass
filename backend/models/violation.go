package models

import (
	"time"

	"gorm.io/gorm"
)

type Violation struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	EmployeeID    *uint          `json:"employee_id"`
	Employee      *Employee      `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	ViolationType string         `json:"violation_type"` // no_helmet, no_vest, etc.
	CameraID      int            `json:"camera_id"`
	Confidence    float64        `json:"confidence"`
	ImageURL      string         `json:"image_url"`
	Status        string         `json:"status"` // new, confirmed, dismissed
	Timestamp     time.Time      `json:"timestamp"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type Statistics struct {
	TotalEmployees  int     `json:"total_employees"`
	ActiveEmployees int     `json:"active_employees"`
	TotalViolations int     `json:"total_violations"`
	ViolationsToday int     `json:"violations_today"`
	ComplianceRate  float64 `json:"compliance_rate"`
}
