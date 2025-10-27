package repository

import (
	"context"
	"time"

	"ppe-detection/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CameraRepository interface {
	Create(ctx context.Context, camera *models.Camera) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Camera, error)
	GetByLocationID(ctx context.Context, locationID uuid.UUID) ([]models.Camera, error)
	GetActive(ctx context.Context) ([]models.Camera, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Camera, error)
	Update(ctx context.Context, camera *models.Camera) error
	UpdateHeartbeat(ctx context.Context, cameraID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type cameraRepo struct {
	db *sqlx.DB
}

func NewCameraRepository(db *sqlx.DB) CameraRepository {
	return &cameraRepo{db: db}
}

func (r *cameraRepo) Create(ctx context.Context, c *models.Camera) error {
	query := `
		INSERT INTO cameras (
			location_id, name, rtsp_url, position, viewing_angle,
			fps, resolution, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		c.LocationID, c.Name, c.RTSPURL, c.Position, c.ViewingAngle,
		c.FPS, c.Resolution, c.IsActive,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *cameraRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Camera, error) {
	var c models.Camera
	query := `SELECT * FROM cameras WHERE id = $1`
	err := r.db.GetContext(ctx, &c, query, id)
	return &c, err
}

func (r *cameraRepo) GetByLocationID(ctx context.Context, locationID uuid.UUID) ([]models.Camera, error) {
	var cameras []models.Camera
	query := `SELECT * FROM cameras WHERE location_id = $1`
	err := r.db.SelectContext(ctx, &cameras, query, locationID)
	return cameras, err
}

func (r *cameraRepo) GetActive(ctx context.Context) ([]models.Camera, error) {
	var cameras []models.Camera
	query := `SELECT * FROM cameras WHERE is_active = true`
	err := r.db.SelectContext(ctx, &cameras, query)
	return cameras, err
}

func (r *cameraRepo) GetAll(ctx context.Context, limit, offset int) ([]models.Camera, error) {
	var cameras []models.Camera
	query := `SELECT * FROM cameras ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &cameras, query, limit, offset)
	return cameras, err
}

func (r *cameraRepo) Update(ctx context.Context, c *models.Camera) error {
	query := `
		UPDATE cameras
		SET name = $2, rtsp_url = $3, position = $4, viewing_angle = $5,
		    fps = $6, resolution = $7, is_active = $8
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.Name, c.RTSPURL, c.Position, c.ViewingAngle,
		c.FPS, c.Resolution, c.IsActive,
	)
	return err
}

func (r *cameraRepo) UpdateHeartbeat(ctx context.Context, cameraID uuid.UUID) error {
	query := `UPDATE cameras SET last_heartbeat = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, cameraID, time.Now())
	return err
}

func (r *cameraRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM cameras WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
