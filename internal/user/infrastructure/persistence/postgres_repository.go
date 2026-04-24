package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/application/ports"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, u *domain.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, email, username, password_hash, status, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		u.ID, u.Email, u.Username, u.PasswordHash, u.Status, u.Version, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, email, username, password_hash, status, version, created_at, updated_at
		FROM users WHERE id = $1`, id)

	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Status, &u.Version, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, domain.ErrNotFound
	}
	return &u, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, email, username, password_hash, status, version, created_at, updated_at
		FROM users WHERE email = $1`, email)

	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Status, &u.Version, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, domain.ErrNotFound
	}
	return &u, nil
}

func (r *PostgresRepository) Update(ctx context.Context, u *domain.User) error {
	u.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, `
		UPDATE users
		SET email=$2, username=$3, password_hash=$4, status=$5, version=$6, updated_at=$7
		WHERE id=$1`,
		u.ID, u.Email, u.Username, u.PasswordHash, u.Status, u.Version, u.UpdatedAt,
	)
	return err
}

var _ ports.UserRepository = (*PostgresRepository)(nil)
