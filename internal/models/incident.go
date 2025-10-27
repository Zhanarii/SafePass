package models

import (
	"time"

	"github.com/google/uuid"
)

type IncidentStatus string

const (
	IncidentOpen          IncidentStatus = "open"
	IncidentInvestigating IncidentStatus = "investigating"
	IncidentResolved      IncidentStatus = "resolved"
	IncidentClosed        IncidentStatus = "closed"
)

type Incident struct {
	ID                       uuid.UUID      `json:"id" db:"id"`
	ViolationID              uuid.UUID      `json:"violation_id" db:"violation_id"`
	CamundaProcessInstanceID *string        `json:"camunda_process_instance_id,omitempty" db:"camunda_process_instance_id"`
	IncidentNumber           string         `json:"incident_number" db:"incident_number"`
	Title                    string         `json:"title" db:"title"`
	Description              string         `json:"description" db:"description"`
	Status                   IncidentStatus `json:"status" db:"status"`
	AssignedTo               *uuid.UUID     `json:"assigned_to,omitempty" db:"assigned_to"`
	RootCause                string         `json:"root_cause,omitempty" db:"root_cause"`
	CorrectiveActions        string         `json:"corrective_actions,omitempty" db:"corrective_actions"`
	DueDate                  *time.Time     `json:"due_date,omitempty" db:"due_date"`
	ResolvedAt               *time.Time     `json:"resolved_at,omitempty" db:"resolved_at"`
	ClosedAt                 *time.Time     `json:"closed_at,omitempty" db:"closed_at"`
	CreatedAt                time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at" db:"updated_at"`
}

type IncidentEvent struct {
	ID         uuid.UUID              `json:"id" db:"id"`
	IncidentID uuid.UUID              `json:"incident_id" db:"incident_id"`
	EventType  string                 `json:"event_type" db:"event_type"`
	UserID     *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	Comment    string                 `json:"comment,omitempty" db:"comment"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`
}
