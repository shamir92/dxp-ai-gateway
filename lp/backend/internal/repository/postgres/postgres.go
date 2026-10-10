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

var ErrNotFound = repository.ErrNotFound

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, u *entity.User) error {
	now := time.Now().UTC()
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = now
	u.UpdatedAt = now
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, email, name, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		u.ID, u.Email, u.Name, u.AvatarURL, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, email, name, avatar_url, created_at, updated_at
		FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, email, name, avatar_url, created_at, updated_at
		FROM users WHERE email = $1`, email)
	return scanUser(row)
}

func (r *UserRepo) Update(ctx context.Context, u *entity.User) error {
	u.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE users SET name = $2, email = $3, avatar_url = $4, updated_at = $5
		WHERE id = $1`,
		u.ID, u.Name, u.Email, u.AvatarURL, u.UpdatedAt)
	return err
}

func scanUser(row pgx.Row) (*entity.User, error) {
	var u entity.User
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

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

// Ensure interface compliance.
var (
	_ repository.UserRepo       = (*UserRepo)(nil)
	_ repository.AuthMethodRepo = (*AuthMethodRepo)(nil)
	_ repository.SessionRepo    = (*SessionRepo)(nil)
	_ repository.SettingsRepo   = (*SettingsRepo)(nil)
	_ repository.MagicCodeRepo  = (*MagicCodeRepo)(nil)
	_ repository.LoginAuditRepo = (*LoginAuditRepo)(nil)
)
