package repository

import (
	"context"

	"ppe-detection/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmployeeID(ctx context.Context, employeeID string) (*models.User, error)
	GetByBadgeNumber(ctx context.Context, badgeNumber string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByRole(ctx context.Context, role models.UserRole) ([]models.User, error)
	GetAvailableSupervisor(ctx context.Context) (*models.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (
			employee_id, first_name, last_name, email, phone,
			role, department, badge_number, photo_url, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		u.EmployeeID, u.FirstName, u.LastName, u.Email, u.Phone,
		u.Role, u.Department, u.BadgeNumber, u.PhotoURL, u.IsActive,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &u, query, id)
	return &u, err
}

func (r *userRepo) GetByEmployeeID(ctx context.Context, employeeID string) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE employee_id = $1`
	err := r.db.GetContext(ctx, &u, query, employeeID)
	return &u, err
}

func (r *userRepo) GetByBadgeNumber(ctx context.Context, badgeNumber string) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE badge_number = $1`
	err := r.db.GetContext(ctx, &u, query, badgeNumber)
	return &u, err
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &u, query, email)
	return &u, err
}

func (r *userRepo) GetByRole(ctx context.Context, role models.UserRole) ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users WHERE role = $1 AND is_active = true`
	err := r.db.SelectContext(ctx, &users, query, role)
	return users, err
}

func (r *userRepo) GetAvailableSupervisor(ctx context.Context) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE role = $1 AND is_active = true ORDER BY RANDOM() LIMIT 1`
	err := r.db.GetContext(ctx, &u, query, models.RoleSupervisor)
	return &u, err
}

func (r *userRepo) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}

func (r *userRepo) Update(ctx context.Context, u *models.User) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3, email = $4, phone = $5,
			role = $6, department = $7, badge_number = $8, photo_url = $9, is_active = $10
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.FirstName, u.LastName, u.Email, u.Phone,
		u.Role, u.Department, u.BadgeNumber, u.PhotoURL, u.IsActive,
	)
	return err
}

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
