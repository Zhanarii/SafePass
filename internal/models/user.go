package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleSupervisor UserRole = "supervisor"
	RoleSecurity   UserRole = "security"
	RoleWorker     UserRole = "worker"
)

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EmployeeID  string    `json:"employee_id" db:"employee_id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Email       string    `json:"email" db:"email"`
	Phone       *string   `json:"phone,omitempty" db:"phone"`
	Role        UserRole  `json:"role" db:"role"`
	Department  *string   `json:"department,omitempty" db:"department"`
	BadgeNumber *string   `json:"badge_number,omitempty" db:"badge_number"`
	PhotoURL    *string   `json:"photo_url,omitempty" db:"photo_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
