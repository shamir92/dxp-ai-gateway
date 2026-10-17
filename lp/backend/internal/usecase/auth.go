package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

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
	users            repository.UserRepo
	methods          repository.AuthMethodRepo
	sessions         repository.SessionRepo
	settings         repository.SettingsRepo
	codes            repository.MagicCodeRepo
	audit            repository.LoginAuditRepo
	newAccounts      repository.NewAccountRepo
	externalAccounts ExternalAccountFactory
	mailer           mail.Mailer
	geo              geoip.Resolver
	now              func() time.Time
	rateWindow       time.Duration
}

type ExternalAccountFactory interface {
	Build(user *entity.User) ([]entity.ExternalSiteAccount, error)
}

type AuthOption func(*AuthUsecase)

func WithNewAccountProvisioning(repo repository.NewAccountRepo, factory ExternalAccountFactory) AuthOption {
	return func(u *AuthUsecase) {
		u.newAccounts = repo
		u.externalAccounts = factory
	}
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
	options ...AuthOption,
) *AuthUsecase {
	if geo == nil {
		geo = geoip.StaticResolver{Label: "Unknown location"}
	}
	u := &AuthUsecase{
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
	for _, option := range options {
		option(u)
	}
	return u
}

func (u *AuthUsecase) createNewAccount(ctx context.Context, user *entity.User, method *entity.AuthMethod) error {
	if u.newAccounts != nil {
		if user.ID == uuid.Nil {
			user.ID = uuid.New()
		}
		var accounts []entity.ExternalSiteAccount
		if u.externalAccounts != nil {
			var err error
			accounts, err = u.externalAccounts.Build(user)
			if err != nil {
				return err
			}
		}
		return u.newAccounts.CreateNewAccount(ctx, user, method, accounts)
	}

	if err := u.users.Create(ctx, user); err != nil {
		return err
	}
	method.UserID = user.ID
	settings := entity.DefaultSettings(user.ID)
	if err := u.settings.Upsert(ctx, &settings); err != nil {
		return err
	}
	return u.methods.Upsert(ctx, method)
}
