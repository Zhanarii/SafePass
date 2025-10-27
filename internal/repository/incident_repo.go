package repository

import (
	"context"
	"encoding/json"
	"time"

	"ppe-detection/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IncidentRepository interface {
	Create(ctx context.Context, incident *models.Incident) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Incident, error)
	GetByViolationID(ctx context.Context, violationID uuid.UUID) (*models.Incident, error)
	GetByStatus(ctx context.Context, status models.IncidentStatus) ([]models.Incident, error)
	GetByAssignedUser(ctx context.Context, userID uuid.UUID) ([]models.Incident, error)
	GetByTimeRange(ctx context.Context, from, to time.Time) ([]models.Incident, error)
	Update(ctx context.Context, incident *models.Incident) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateEvent(ctx context.Context, event *models.IncidentEvent) error
	GetEventsByIncidentID(ctx context.Context, incidentID uuid.UUID) ([]models.IncidentEvent, error)
}

type incidentRepo struct {
	db *sqlx.DB
}

func NewIncidentRepository(db *sqlx.DB) IncidentRepository {
	return &incidentRepo{db: db}
}

func (r *incidentRepo) Create(ctx context.Context, i *models.Incident) error {
	query := `
		INSERT INTO incidents (
			violation_id, camunda_process_instance_id, incident_number,
			title, description, status, assigned_to, due_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		i.ViolationID, i.CamundaProcessInstanceID, i.IncidentNumber,
		i.Title, i.Description, i.Status, i.AssignedTo, i.DueDate,
	).Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
}

func (r *incidentRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Incident, error) {
	var i models.Incident
	query := `SELECT * FROM incidents WHERE id = $1`
	err := r.db.GetContext(ctx, &i, query, id)
	return &i, err
}

func (r *incidentRepo) GetByViolationID(ctx context.Context, violationID uuid.UUID) (*models.Incident, error) {
	var i models.Incident
	query := `SELECT * FROM incidents WHERE violation_id = $1`
	err := r.db.GetContext(ctx, &i, query, violationID)
	return &i, err
}

func (r *incidentRepo) GetByStatus(ctx context.Context, status models.IncidentStatus) ([]models.Incident, error) {
	var incidents []models.Incident
	query := `SELECT * FROM incidents WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &incidents, query, status)
	return incidents, err
}

func (r *incidentRepo) GetByAssignedUser(ctx context.Context, userID uuid.UUID) ([]models.Incident, error) {
	var incidents []models.Incident
	query := `SELECT * FROM incidents WHERE assigned_to = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &incidents, query, userID)
	return incidents, err
}

func (r *incidentRepo) GetByTimeRange(ctx context.Context, from, to time.Time) ([]models.Incident, error) {
	var incidents []models.Incident
	query := `SELECT * FROM incidents WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &incidents, query, from, to)
	return incidents, err
}

func (r *incidentRepo) Update(ctx context.Context, i *models.Incident) error {
	query := `
		UPDATE incidents
		SET camunda_process_instance_id = $2, status = $3, assigned_to = $4,
		    root_cause = $5, corrective_actions = $6, due_date = $7,
		    resolved_at = $8, closed_at = $9
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		i.ID, i.CamundaProcessInstanceID, i.Status, i.AssignedTo,
		i.RootCause, i.CorrectiveActions, i.DueDate,
		i.ResolvedAt, i.ClosedAt,
	)
	return err
}

func (r *incidentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM incidents WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *incidentRepo) CreateEvent(ctx context.Context, e *models.IncidentEvent) error {
	query := `
		INSERT INTO incident_events (
			incident_id, event_type, user_id, comment, metadata
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	metadataJSON, _ := json.Marshal(e.Metadata)
	return r.db.QueryRowContext(ctx, query,
		e.IncidentID, e.EventType, e.UserID, e.Comment, metadataJSON,
	).Scan(&e.ID, &e.CreatedAt)
}

func (r *incidentRepo) GetEventsByIncidentID(ctx context.Context, incidentID uuid.UUID) ([]models.IncidentEvent, error) {
	var events []models.IncidentEvent
	query := `SELECT * FROM incident_events WHERE incident_id = $1 ORDER BY created_at ASC`
	err := r.db.SelectContext(ctx, &events, query, incidentID)
	return events, err
}
