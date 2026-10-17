package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
)

type LoginAuditRepo struct {
	pool *pgxpool.Pool
}

func NewLoginAuditRepo(pool *pgxpool.Pool) *LoginAuditRepo {
	return &LoginAuditRepo{pool: pool}
}

func (r *LoginAuditRepo) Insert(ctx context.Context, a *entity.LoginAudit) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now().UTC()
	}
	var ip any
	if a.IP != "" {
		ip = a.IP
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO login_audit (id, user_id, session_id, ip, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		a.ID, a.UserID, a.SessionID, ip, a.UserAgent, a.CreatedAt)
	return err
}

func (r *LoginAuditRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]entity.LoginAudit, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, session_id, host(ip), user_agent, created_at
		FROM login_audit WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.LoginAudit
	for rows.Next() {
		var a entity.LoginAudit
		var ip *string
		if err := rows.Scan(&a.ID, &a.UserID, &a.SessionID, &ip, &a.UserAgent, &a.CreatedAt); err != nil {
			return nil, err
		}
		if ip != nil {
			a.IP = *ip
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

var _ repository.LoginAuditRepo = (*LoginAuditRepo)(nil)
