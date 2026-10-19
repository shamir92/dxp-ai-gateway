package repository

import (
	"context"
	"errors"
	"time"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("repository: not found")
var ErrLeaseLost = errors.New("repository: external account lease lost")

type UserRepo interface {
	Create(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, u *entity.User) error
}

type AuthMethodRepo interface {
	Upsert(ctx context.Context, m *entity.AuthMethod) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]entity.AuthMethod, error)
	GetByProvider(ctx context.Context, provider entity.Provider, providerUserID string) (*entity.AuthMethod, error)
	Delete(ctx context.Context, userID uuid.UUID, provider entity.Provider) error
	CountByUser(ctx context.Context, userID uuid.UUID) (int, error)
}

type SessionRepo interface {
	Create(ctx context.Context, s *entity.Session) error
	GetByTokenHash(ctx context.Context, hash string) (*entity.Session, error)
	Touch(ctx context.Context, id uuid.UUID, at time.Time) error
	Revoke(ctx context.Context, id uuid.UUID, at time.Time) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID, at time.Time) error
	ListActive(ctx context.Context, userID uuid.UUID) ([]entity.Session, error)
}

type SettingsRepo interface {
	Get(ctx context.Context, userID uuid.UUID) (*entity.UserSettings, error)
	Upsert(ctx context.Context, s *entity.UserSettings) error
}

type MagicCodeRepo interface {
	Insert(ctx context.Context, email, name, codeHash string, purpose entity.MagicPurpose, expiresAt time.Time) error
	GetLatestActive(ctx context.Context, email string, purpose entity.MagicPurpose, now time.Time) (id uuid.UUID, name, codeHash string, attempts int, err error)
	InvalidateForEmail(ctx context.Context, email string, purpose entity.MagicPurpose) error
	IncrementAttempts(ctx context.Context, id uuid.UUID) (attempts int, err error)
	MarkConsumed(ctx context.Context, id uuid.UUID, at time.Time) error
}

type LoginAuditRepo interface {
	Insert(ctx context.Context, a *entity.LoginAudit) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]entity.LoginAudit, error)
}

type NewAccountRepo interface {
	CreateNewAccount(ctx context.Context, user *entity.User, method *entity.AuthMethod, accounts []entity.ExternalSiteAccount) error
}

type ExternalSiteAccountRepo interface {
	ClaimDue(ctx context.Context, now, leaseExpiresAt time.Time, limit int) ([]entity.ExternalSiteAccount, error)
	MarkProvisioned(ctx context.Context, id, leaseToken uuid.UUID, remoteUserID string, at time.Time) error
	MarkRetry(ctx context.Context, id, leaseToken uuid.UUID, nextAttemptAt time.Time, summary string, at time.Time) error
	MarkFailed(ctx context.Context, id, leaseToken uuid.UUID, summary string, at time.Time) error
	// GetByUserAndSite reads one user's provisioned account so a caller can obtain
	// the remote username and the encrypted site password.
	GetByUserAndSite(ctx context.Context, userID uuid.UUID, site entity.ExternalSite) (*entity.ExternalSiteAccount, error)
	// EnsurePending inserts the given per-site rows for a user that has none,
	// leaving existing rows untouched. It backfills users who predate the
	// provisioning flow so the worker can create their remote accounts.
	EnsurePending(ctx context.Context, accounts []entity.ExternalSiteAccount) error
	// Requeue marks a row provisioned-but-unusable for provisioning again, which is
	// how a remote account that no longer exists gets recreated. It reports whether
	// a row was actually requeued and only touches rows the worker believes are
	// provisioned, up to maxAttempts.
	Requeue(ctx context.Context, userID uuid.UUID, site entity.ExternalSite, maxAttempts int) (bool, error)
}

// AccessTokenRepo stores the encrypted personal access token ("nap_…") the
// gateway reuses for per-user calls to an external site. It never exposes a
// plaintext token.
type AccessTokenRepo interface {
	GetAccessToken(ctx context.Context, userID uuid.UUID, site entity.ExternalSite) (*entity.StoredAccessToken, error)
	SaveAccessToken(ctx context.Context, token *entity.StoredAccessToken) error
	ClearAccessToken(ctx context.Context, userID uuid.UUID, site entity.ExternalSite) error
}

type Repos struct {
	Users        UserRepo
	AuthMethods  AuthMethodRepo
	Sessions     SessionRepo
	Settings     SettingsRepo
	MagicCodes   MagicCodeRepo
	LoginAudit   LoginAuditRepo
	AccessTokens AccessTokenRepo
}
