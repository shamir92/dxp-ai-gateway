package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
)

type AuthMethodRepo struct {
	pool *pgxpool.Pool
}

func NewAuthMethodRepo(pool *pgxpool.Pool) *AuthMethodRepo {
	return &AuthMethodRepo{pool: pool}
}

func (r *AuthMethodRepo) Upsert(ctx context.Context, m *entity.AuthMethod) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO auth_methods (id, user_id, provider, provider_user_id, email, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (provider, provider_user_id) DO UPDATE
		SET user_id = EXCLUDED.user_id, email = EXCLUDED.email`,
		m.ID, m.UserID, m.Provider, m.ProviderUserID, m.Email, m.CreatedAt)
	return err
}

func (r *AuthMethodRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]entity.AuthMethod, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, provider, provider_user_id, email, created_at
		FROM auth_methods WHERE user_id = $1 ORDER BY provider`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.AuthMethod
	for rows.Next() {
		var m entity.AuthMethod
		if err := rows.Scan(&m.ID, &m.UserID, &m.Provider, &m.ProviderUserID, &m.Email, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *AuthMethodRepo) GetByProvider(ctx context.Context, provider entity.Provider, providerUserID string) (*entity.AuthMethod, error) {
	var m entity.AuthMethod
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, provider, provider_user_id, email, created_at
		FROM auth_methods WHERE provider = $1 AND provider_user_id = $2`,
		provider, providerUserID).
		Scan(&m.ID, &m.UserID, &m.Provider, &m.ProviderUserID, &m.Email, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *AuthMethodRepo) Delete(ctx context.Context, userID uuid.UUID, provider entity.Provider) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM auth_methods WHERE user_id = $1 AND provider = $2`, userID, provider)
	return err
}

func (r *AuthMethodRepo) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM auth_methods WHERE user_id = $1`, userID).Scan(&n)
	return n, err
}

var _ repository.AuthMethodRepo = (*AuthMethodRepo)(nil)
