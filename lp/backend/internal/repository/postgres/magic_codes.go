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

type MagicCodeRepo struct {
	pool *pgxpool.Pool
}

func NewMagicCodeRepo(pool *pgxpool.Pool) *MagicCodeRepo {
	return &MagicCodeRepo{pool: pool}
}

func (r *MagicCodeRepo) Insert(ctx context.Context, email, name, codeHash string, purpose entity.MagicPurpose, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO magic_codes (email, name, code_hash, purpose, expires_at)
		VALUES ($1, $2, $3, $4, $5)`, email, name, codeHash, purpose, expiresAt)
	return err
}

func (r *MagicCodeRepo) GetLatestActive(ctx context.Context, email string, purpose entity.MagicPurpose, now time.Time) (uuid.UUID, string, string, int, error) {
	var id uuid.UUID
	var name, hash string
	var attempts int
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, code_hash, attempts FROM magic_codes
		WHERE email = $1 AND purpose = $2 AND consumed_at IS NULL AND expires_at > $3
		ORDER BY created_at DESC LIMIT 1`, email, purpose, now).
		Scan(&id, &name, &hash, &attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, "", "", 0, ErrNotFound
	}
	return id, name, hash, attempts, err
}

func (r *MagicCodeRepo) InvalidateForEmail(ctx context.Context, email string, purpose entity.MagicPurpose) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE magic_codes SET consumed_at = now()
		WHERE email = $1 AND purpose = $2 AND consumed_at IS NULL`, email, purpose)
	return err
}

func (r *MagicCodeRepo) IncrementAttempts(ctx context.Context, id uuid.UUID) (int, error) {
	var attempts int
	err := r.pool.QueryRow(ctx, `
		UPDATE magic_codes SET attempts = attempts + 1
		WHERE id = $1 RETURNING attempts`, id).Scan(&attempts)
	return attempts, err
}

func (r *MagicCodeRepo) MarkConsumed(ctx context.Context, id uuid.UUID, at time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE magic_codes SET consumed_at = $2 WHERE id = $1`, id, at)
	return err
}

var _ repository.MagicCodeRepo = (*MagicCodeRepo)(nil)
