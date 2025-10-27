package repository

import (
	"context"
	"time"

	"ppe-detection/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ViolationRepository interface {
	Create(ctx context.Context, violation *models.Violation) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Violation, error)
	GetByDetectionID(ctx context.Context, detectionID uuid.UUID) (*models.Violation, error)
	GetByTimeRange(ctx context.Context, from, to time.Time) ([]models.Violation, error)
	GetByUserAndTimeRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Violation, error)
	GetByCameraAndTimeRange(ctx context.Context, cameraID uuid.UUID, from, to time.Time) ([]models.Violation, error)
	GetUnacknowledged(ctx context.Context, limit int) ([]models.Violation, error)
	Update(ctx context.Context, violation *models.Violation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type violationRepo struct {
	db *sqlx.DB
}

func NewViolationRepository(db *sqlx.DB) ViolationRepository {
	return &violationRepo{db: db}
}

func (r *violationRepo) Create(ctx context.Context, v *models.Violation) error {
	query := `
		INSERT INTO violations (
			detection_id, user_id, camera_id, access_zone_id, violation_type,
			severity, description, snapshot_url, video_url, access_decision
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		v.DetectionID, v.UserID, v.CameraID, v.AccessZoneID, v.ViolationType,
		v.Severity, v.Description, v.SnapshotURL, v.VideoURL, v.AccessDecision,
	).Scan(&v.ID, &v.CreatedAt)
}

func (r *violationRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Violation, error) {
	var v models.Violation
	query := `SELECT * FROM violations WHERE id = $1`
	err := r.db.GetContext(ctx, &v, query, id)
	return &v, err
}

func (r *violationRepo) GetByDetectionID(ctx context.Context, detectionID uuid.UUID) (*models.Violation, error) {
	var v models.Violation
	query := `SELECT * FROM violations WHERE detection_id = $1`
	err := r.db.GetContext(ctx, &v, query, detectionID)
	return &v, err
}

func (r *violationRepo) GetByTimeRange(ctx context.Context, from, to time.Time) ([]models.Violation, error) {
	var violations []models.Violation
	query := `SELECT * FROM violations WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &violations, query, from, to)
	return violations, err
}

func (r *violationRepo) GetByUserAndTimeRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Violation, error) {
	var violations []models.Violation
	query := `SELECT * FROM violations WHERE user_id = $1 AND created_at BETWEEN $2 AND $3 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &violations, query, userID, from, to)
	return violations, err
}

func (r *violationRepo) GetByCameraAndTimeRange(ctx context.Context, cameraID uuid.UUID, from, to time.Time) ([]models.Violation, error) {
	var violations []models.Violation
	query := `SELECT * FROM violations WHERE camera_id = $1 AND created_at BETWEEN $2 AND $3 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &violations, query, cameraID, from, to)
	return violations, err
}

func (r *violationRepo) GetUnacknowledged(ctx context.Context, limit int) ([]models.Violation, error) {
	var violations []models.Violation
	query := `SELECT * FROM violations WHERE acknowledged_at IS NULL ORDER BY created_at DESC LIMIT $1`
	err := r.db.SelectContext(ctx, &violations, query, limit)
	return violations, err
}

func (r *violationRepo) Update(ctx context.Context, v *models.Violation) error {
	query := `UPDATE violations SET acknowledged_by = $2, acknowledged_at = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, v.ID, v.AcknowledgedBy, v.AcknowledgedAt)
	return err
}

func (r *violationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM violations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
