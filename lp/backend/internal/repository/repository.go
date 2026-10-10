package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
)

var ErrNotFound = errors.New("repository: not found")

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

type Repos struct {
	Users      UserRepo
	AuthMethods AuthMethodRepo
	Sessions   SessionRepo
	Settings   SettingsRepo
	MagicCodes MagicCodeRepo
	LoginAudit LoginAuditRepo
}
