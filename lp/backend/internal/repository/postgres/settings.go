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

type SettingsRepo struct {
	pool *pgxpool.Pool
}

func NewSettingsRepo(pool *pgxpool.Pool) *SettingsRepo {
	return &SettingsRepo{pool: pool}
}

func (r *SettingsRepo) Get(ctx context.Context, userID uuid.UUID) (*entity.UserSettings, error) {
	var s entity.UserSettings
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, default_region, default_model, notify_balance, notify_usage,
		       notify_security, notify_product, stream_default, store_request_logs, login_alerts, updated_at
		FROM user_settings WHERE user_id = $1`, userID).
		Scan(&s.UserID, &s.DefaultRegion, &s.DefaultModel, &s.NotifyBalance, &s.NotifyUsage,
			&s.NotifySecurity, &s.NotifyProduct, &s.StreamDefault, &s.StoreRequestLogs, &s.LoginAlerts, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SettingsRepo) Upsert(ctx context.Context, s *entity.UserSettings) error {
	s.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_settings (
			user_id, default_region, default_model, notify_balance, notify_usage,
			notify_security, notify_product, stream_default, store_request_logs, login_alerts, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (user_id) DO UPDATE SET
			default_region = EXCLUDED.default_region,
			default_model = EXCLUDED.default_model,
			notify_balance = EXCLUDED.notify_balance,
			notify_usage = EXCLUDED.notify_usage,
			notify_security = EXCLUDED.notify_security,
			notify_product = EXCLUDED.notify_product,
			stream_default = EXCLUDED.stream_default,
			store_request_logs = EXCLUDED.store_request_logs,
			login_alerts = EXCLUDED.login_alerts,
			updated_at = EXCLUDED.updated_at`,
		s.UserID, s.DefaultRegion, s.DefaultModel, s.NotifyBalance, s.NotifyUsage,
		s.NotifySecurity, s.NotifyProduct, s.StreamDefault, s.StoreRequestLogs, s.LoginAlerts, s.UpdatedAt)
	return err
}

var _ repository.SettingsRepo = (*SettingsRepo)(nil)
