package repository

import (
	"context"
	"encoding/json"

	"ppe-detection/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccessZoneRepository interface {
	Create(ctx context.Context, zone *models.AccessZone) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.AccessZone, error)
	GetByLocationID(ctx context.Context, locationID uuid.UUID) ([]models.AccessZone, error)
	GetByCameraID(ctx context.Context, cameraID uuid.UUID) (*models.AccessZone, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.AccessZone, error)
	Update(ctx context.Context, zone *models.AccessZone) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type accessZoneRepo struct {
	db *sqlx.DB
}

func NewAccessZoneRepository(db *sqlx.DB) AccessZoneRepository {
	return &accessZoneRepo{db: db}
}

func (r *accessZoneRepo) Create(ctx context.Context, z *models.AccessZone) error {
	query := `
		INSERT INTO access_zones (
			location_id, camera_id, name, description, required_ppe,
			danger_level, access_rules, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	rulesJSON, _ := json.Marshal(z.AccessRules)
	return r.db.QueryRowContext(ctx, query,
		z.LocationID, z.CameraID, z.Name, z.Description, z.RequiredPPE,
		z.DangerLevel, rulesJSON, z.IsActive,
	).Scan(&z.ID, &z.CreatedAt, &z.UpdatedAt)
}

func (r *accessZoneRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.AccessZone, error) {
	var z models.AccessZone
	query := `SELECT * FROM access_zones WHERE id = $1`
	err := r.db.GetContext(ctx, &z, query, id)
	return &z, err
}

func (r *accessZoneRepo) GetByLocationID(ctx context.Context, locationID uuid.UUID) ([]models.AccessZone, error) {
	var zones []models.AccessZone
	query := `SELECT * FROM access_zones WHERE location_id = $1 AND is_active = true`
	err := r.db.SelectContext(ctx, &zones, query, locationID)
	return zones, err
}

func (r *accessZoneRepo) GetByCameraID(ctx context.Context, cameraID uuid.UUID) (*models.AccessZone, error) {
	var z models.AccessZone
	query := `SELECT * FROM access_zones WHERE camera_id = $1`
	err := r.db.GetContext(ctx, &z, query, cameraID)
	return &z, err
}

func (r *accessZoneRepo) GetAll(ctx context.Context, limit, offset int) ([]models.AccessZone, error) {
	var zones []models.AccessZone
	query := `SELECT * FROM access_zones ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &zones, query, limit, offset)
	return zones, err
}

func (r *accessZoneRepo) Update(ctx context.Context, z *models.AccessZone) error {
	query := `
		UPDATE access_zones
		SET camera_id = $2, name = $3, description = $4, required_ppe = $5,
		    danger_level = $6, access_rules = $7, is_active = $8
		WHERE id = $1`
	rulesJSON, _ := json.Marshal(z.AccessRules)
	_, err := r.db.ExecContext(ctx, query,
		z.ID, z.CameraID, z.Name, z.Description, z.RequiredPPE,
		z.DangerLevel, rulesJSON, z.IsActive,
	)
	return err
}

func (r *accessZoneRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM access_zones WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
