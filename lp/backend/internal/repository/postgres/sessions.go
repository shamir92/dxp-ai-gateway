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

type SessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{pool: pool}
}

func (r *SessionRepo) Create(ctx context.Context, s *entity.Session) error {
	now := time.Now().UTC()
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	if s.LastSeenAt.IsZero() {
		s.LastSeenAt = now
	}
	var ip any
	if s.IP != "" {
		ip = s.IP
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sessions (id, token_hash, user_id, user_agent, ip, location, created_at, last_seen_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		s.ID, s.TokenHash, s.UserID, s.UserAgent, ip, s.Location, s.CreatedAt, s.LastSeenAt)
	return err
}

func (r *SessionRepo) GetByTokenHash(ctx context.Context, hash string) (*entity.Session, error) {
	var s entity.Session
	var ip *string
	err := r.pool.QueryRow(ctx, `
		SELECT id, token_hash, user_id, user_agent, host(ip), location, created_at, last_seen_at, revoked_at
		FROM sessions WHERE token_hash = $1`, hash).
		Scan(&s.ID, &s.TokenHash, &s.UserID, &s.UserAgent, &ip, &s.Location, &s.CreatedAt, &s.LastSeenAt, &s.RevokedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if ip != nil {
		s.IP = *ip
	}
	return &s, nil
}

func (r *SessionRepo) Touch(ctx context.Context, id uuid.UUID, at time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE sessions SET last_seen_at = $2 WHERE id = $1`, id, at)
	return err
}

func (r *SessionRepo) Revoke(ctx context.Context, id uuid.UUID, at time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE sessions SET revoked_at = $2 WHERE id = $1 AND revoked_at IS NULL`, id, at)
	return err
}

func (r *SessionRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID, at time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE sessions SET revoked_at = $2 WHERE user_id = $1 AND revoked_at IS NULL`, userID, at)
	return err
}

func (r *SessionRepo) ListActive(ctx context.Context, userID uuid.UUID) ([]entity.Session, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, token_hash, user_id, user_agent, host(ip), location, created_at, last_seen_at, revoked_at
		FROM sessions WHERE user_id = $1 AND revoked_at IS NULL
		ORDER BY last_seen_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.Session
	for rows.Next() {
		var s entity.Session
		var ip *string
		if err := rows.Scan(&s.ID, &s.TokenHash, &s.UserID, &s.UserAgent, &ip, &s.Location, &s.CreatedAt, &s.LastSeenAt, &s.RevokedAt); err != nil {
			return nil, err
		}
		if ip != nil {
			s.IP = *ip
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

var _ repository.SessionRepo = (*SessionRepo)(nil)
