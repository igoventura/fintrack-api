package postgres

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-api/domain"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, supabase_id, name, email, created_at, updated_at, deactivated_at FROM users WHERE id = $1 AND deactivated_at IS NULL`
	var u domain.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.SupabaseID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeactivatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, supabase_id, name, email, created_at, updated_at, deactivated_at FROM users WHERE email = $1 AND deactivated_at IS NULL`
	var u domain.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.SupabaseID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeactivatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetBySupabaseID(ctx context.Context, supabaseID string) (*domain.User, error) {
	query := `SELECT id, supabase_id, name, email, created_at, updated_at, deactivated_at FROM users WHERE supabase_id = $1 AND deactivated_at IS NULL`
	var u domain.User
	err := r.db.Pool.QueryRow(ctx, query, supabaseID).Scan(
		&u.ID, &u.SupabaseID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeactivatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by supabase id: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (supabase_id, name, email)
			  VALUES ($1, $2, $3)
			  RETURNING id, created_at, updated_at`
	row := r.db.Pool.QueryRow(ctx, query, u.SupabaseID, u.Name, u.Email)
	if err := row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	query := `UPDATE users SET name = $2, email = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING updated_at`
	row := r.db.Pool.QueryRow(ctx, query, u.ID, u.Name, u.Email)
	if err := row.Scan(&u.UpdatedAt); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET deactivated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *UserRepository) AddUserToTenant(ctx context.Context, userID, tenantID string) error {
	query := `INSERT INTO users_tenants (user_id, tenant_id) VALUES ($1, $2)`
	_, err := r.db.Pool.Exec(ctx, query, userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to add user to tenant: %w", err)
	}
	return nil
}

func (r *UserRepository) RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error {
	query := `UPDATE users_tenants SET deactivated_at = CURRENT_TIMESTAMP WHERE user_id = $1 AND tenant_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to remove user from tenant: %w", err)
	}
	return nil
}

func (r *UserRepository) ListUserTenants(ctx context.Context, userID string) ([]domain.Tenant, error) {
	query := `SELECT t.id, t.name, t.created_at, t.updated_at, t.deactivated_at
			  FROM users_tenants ut
			  JOIN tenants t ON ut.tenant_id = t.id
			  WHERE user_id = $1 AND ut.deactivated_at IS NULL AND t.deactivated_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user tenants: %w", err)
	}
	defer rows.Close()

	var tenants []domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt, &t.DeactivatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user tenant: %w", err)
		}
		tenants = append(tenants, t)
	}
	return tenants, nil
}
