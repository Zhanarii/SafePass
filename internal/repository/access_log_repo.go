package repository

import (
	"context"
	"time"

	"ppe-detection/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccessLogRepository interface {
	Create(ctx context.Context, log *models.AccessLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.AccessLog, error)
	GetByUserAndTimeRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.AccessLog, error)
	GetByCameraAndTimeRange(ctx context.Context, cameraID uuid.UUID, from, to time.Time) ([]models.AccessLog, error)
	GetByAccessZoneAndTimeRange(ctx context.Context, zoneID uuid.UUID, from, to time.Time) ([]models.AccessLog, error)
	Update(ctx context.Context, log *models.AccessLog) error
}

type accessLogRepo struct {
	db *sqlx.DB
}

func NewAccessLogRepository(db *sqlx.DB) AccessLogRepository {
	return &accessLogRepo{db: db}
}

func (r *accessLogRepo) Create(ctx context.Context, log *models.AccessLog) error {
	query := `
		INSERT INTO access_logs (
			user_id, camera_id, access_zone_id, detection_id,
			decision, badge_scanned, override_by, override_reason
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, timestamp`
	return r.db.QueryRowContext(ctx, query,
		log.UserID, log.CameraID, log.AccessZoneID, log.DetectionID,
		log.Decision, log.BadgeScanned, log.OverrideBy, log.OverrideReason,
	).Scan(&log.ID, &log.Timestamp)
}

func (r *accessLogRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.AccessLog, error) {
	var log models.AccessLog
	query := `SELECT * FROM access_logs WHERE id = $1`
	err := r.db.GetContext(ctx, &log, query, id)
	return &log, err
}

func (r *accessLogRepo) GetByUserAndTimeRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.AccessLog, error) {
	var logs []models.AccessLog
	query := `SELECT * FROM access_logs WHERE user_id = $1 AND timestamp BETWEEN $2 AND $3 ORDER BY timestamp DESC`
	err := r.db.SelectContext(ctx, &logs, query, userID, from, to)
	return logs, err
}

func (r *accessLogRepo) GetByCameraAndTimeRange(ctx context.Context, cameraID uuid.UUID, from, to time.Time) ([]models.AccessLog, error) {
	var logs []models.AccessLog
	query := `SELECT * FROM access_logs WHERE camera_id = $1 AND timestamp BETWEEN $2 AND $3 ORDER BY timestamp DESC`
	err := r.db.SelectContext(ctx, &logs, query, cameraID, from, to)
	return logs, err
}

func (r *accessLogRepo) GetByAccessZoneAndTimeRange(ctx context.Context, zoneID uuid.UUID, from, to time.Time) ([]models.AccessLog, error) {
	var logs []models.AccessLog
	query := `SELECT * FROM access_logs WHERE access_zone_id = $1 AND timestamp BETWEEN $2 AND $3 ORDER BY timestamp DESC`
	err := r.db.SelectContext(ctx, &logs, query, zoneID, from, to)
	return logs, err
}

func (r *accessLogRepo) Update(ctx context.Context, log *models.AccessLog) error {
	query := `UPDATE access_logs SET decision = $2, override_by = $3, override_reason = $4 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.Decision, log.OverrideBy, log.OverrideReason)
	return err
}
