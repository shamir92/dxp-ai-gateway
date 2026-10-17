package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/geoip"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/mail"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCodeExpired        = errors.New("code expired or not found")
	ErrTooManyAttempts    = errors.New("too many attempts")
	ErrRateLimited        = errors.New("rate limited")
	ErrLastAuthMethod     = errors.New("cannot remove last auth method")
	ErrConflict           = errors.New("conflict")
	ErrValidation         = errors.New("validation")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrNotFound           = errors.New("not found")
	ErrNoAccount          = errors.New("no account for this email")
	ErrAccountExists      = errors.New("account already exists")
)

type SessionMeta struct {
	UserAgent string
	IP        string
}

type AuthResult struct {
	User    *entity.User
	Session *entity.Session
	Token   string // raw token for cookie
}

type AuthUsecase struct {
	users      repository.UserRepo
	methods    repository.AuthMethodRepo
	sessions   repository.SessionRepo
	settings   repository.SettingsRepo
	codes      repository.MagicCodeRepo
	audit      repository.LoginAuditRepo
	mailer     mail.Mailer
	geo        geoip.Resolver
	now        func() time.Time
	rateWindow time.Duration
}

func NewAuthUsecase(
	users repository.UserRepo,
	methods repository.AuthMethodRepo,
	sessions repository.SessionRepo,
	settings repository.SettingsRepo,
	codes repository.MagicCodeRepo,
	audit repository.LoginAuditRepo,
	mailer mail.Mailer,
	geo geoip.Resolver,
) *AuthUsecase {
	if geo == nil {
		geo = geoip.StaticResolver{Label: "Unknown location"}
	}
	return &AuthUsecase{
		users:      users,
		methods:    methods,
		sessions:   sessions,
		settings:   settings,
		codes:      codes,
		audit:      audit,
		mailer:     mailer,
		geo:        geo,
		now:        func() time.Time { return time.Now().UTC() },
		rateWindow: time.Hour,
	}
}
