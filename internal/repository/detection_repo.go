package repository

import (
	"context"
	"encoding/json"
	"time"

	"ppe-detection/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DetectionRepository interface {
	Create(ctx context.Context, detection *models.Detection) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Detection, error)
	GetByCameraAndTimeRange(ctx context.Context, cameraID uuid.UUID, from, to time.Time, status *models.DetectionStatus) ([]models.Detection, error)
	GetByUserAndTimeRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Detection, error)
	GetRecentByCameraID(ctx context.Context, cameraID uuid.UUID, limit int) ([]models.Detection, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type detectionRepo struct {
	db *sqlx.DB
}

func NewDetectionRepository(db *sqlx.DB) DetectionRepository {
	return &detectionRepo{db: db}
}

func (r *detectionRepo) Create(ctx context.Context, d *models.Detection) error {
	query := `
		INSERT INTO detections (
			camera_id, user_id, access_zone_id, timestamp, frame_url,
			detected_ppe, missing_ppe, confidence_scores, bounding_boxes,
			status, face_embedding, processing_time_ms, model_version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`

	confidenceJSON, _ := json.Marshal(d.ConfidenceScores)
	bboxJSON, _ := json.Marshal(d.BoundingBoxes)

	return r.db.QueryRowContext(ctx, query,
		d.CameraID, d.UserID, d.AccessZoneID, d.Timestamp, d.FrameURL,
		d.DetectedPPE, d.MissingPPE, confidenceJSON, bboxJSON,
		d.Status, d.FaceEmbedding, d.ProcessingTimeMS, d.ModelVersion,
	).Scan(&d.ID)
}

func (r *detectionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Detection, error) {
	var d models.Detection
	query := `SELECT * FROM detections WHERE id = $1`
	err := r.db.GetContext(ctx, &d, query, id)
	return &d, err
}

func (r *detectionRepo) GetByCameraAndTimeRange(
	ctx context.Context,
	cameraID uuid.UUID,
	from, to time.Time,
	status *models.DetectionStatus,
) ([]models.Detection, error) {
	var detections []models.Detection
	query := `SELECT * FROM detections WHERE camera_id = $1 AND timestamp BETWEEN $2 AND $3`
	args := []interface{}{cameraID, from, to}
	if status != nil {
		query += ` AND status = $4`
		args = append(args, *status)
	}
	query += ` ORDER BY timestamp DESC`
	err := r.db.SelectContext(ctx, &detections, query, args...)
	return detections, err
}

func (r *detectionRepo) GetByUserAndTimeRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Detection, error) {
	var detections []models.Detection
	query := `SELECT * FROM detections WHERE user_id = $1 AND timestamp BETWEEN $2 AND $3 ORDER BY timestamp DESC`
	err := r.db.SelectContext(ctx, &detections, query, userID, from, to)
	return detections, err
}

func (r *detectionRepo) GetRecentByCameraID(ctx context.Context, cameraID uuid.UUID, limit int) ([]models.Detection, error) {
	var detections []models.Detection
	query := `SELECT * FROM detections WHERE camera_id = $1 ORDER BY timestamp DESC LIMIT $2`
	err := r.db.SelectContext(ctx, &detections, query, cameraID, limit)
	return detections, err
}

func (r *detectionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM detections WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
