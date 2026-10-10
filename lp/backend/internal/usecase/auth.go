package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/geoip"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/mail"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/ua"
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

const (
	magicCodeTTL      = 10 * time.Minute
	magicMaxAttempts  = 5
	magicResendWindow = 60 * time.Second
	magicRatePerHour  = 5
	magicExpiresMin   = 10
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

// lastRequestAt is in-memory rate state for magic requests (per process).
type rateEntry struct {
	times []time.Time
}

var magicRate = map[string]*rateEntry{}

func (u *AuthUsecase) RequestMagicCode(ctx context.Context, email, purpose, name string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	p := entity.MagicPurpose(purpose)
	if p != entity.PurposeLogin && p != entity.PurposeRegister {
		return fmt.Errorf("%w: purpose", ErrValidation)
	}
	if email == "" || !strings.Contains(email, "@") {
		return fmt.Errorf("%w: email", ErrValidation)
	}
	if p == entity.PurposeRegister && strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name required for register", ErrValidation)
	}

	// Login requires an existing account; register requires a new one.
	if p == entity.PurposeLogin {
		if _, err := u.users.GetByEmail(ctx, email); err != nil {
			return ErrNoAccount
		}
	} else if p == entity.PurposeRegister {
		if _, err := u.users.GetByEmail(ctx, email); err == nil {
			return ErrAccountExists
		}
	}

	if err := u.allowMagicRequest(email + ":" + purpose); err != nil {
		return err
	}

	code, err := generateNumericCode()
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := u.codes.InvalidateForEmail(ctx, email, p); err != nil {
		return err
	}
	if err := u.codes.Insert(ctx, email, strings.TrimSpace(name), string(hash), p, u.now().Add(magicCodeTTL)); err != nil {
		return err
	}
	return u.mailer.SendMagicCode(ctx, mail.MagicCodeMessage{
		To:               email,
		Code:             code,
		Purpose:          string(p),
		ExpiresInMinutes: magicExpiresMin,
	})
}

func (u *AuthUsecase) ResendMagicCode(ctx context.Context, email, purpose string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	p := entity.MagicPurpose(purpose)
	if p != entity.PurposeLogin && p != entity.PurposeRegister {
		return fmt.Errorf("%w: purpose", ErrValidation)
	}
	if email == "" || !strings.Contains(email, "@") {
		return fmt.Errorf("%w: email", ErrValidation)
	}
	key := email + ":" + purpose
	if !u.allowResendCooldown(key) {
		return ErrRateLimited
	}
	if p == entity.PurposeLogin {
		if _, err := u.users.GetByEmail(ctx, email); err != nil {
			return ErrNoAccount
		}
	}
	if err := u.allowMagicRequest(key); err != nil {
		return err
	}
	code, err := generateNumericCode()
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := u.codes.InvalidateForEmail(ctx, email, p); err != nil {
		return err
	}
	if err := u.codes.Insert(ctx, email, "", string(hash), p, u.now().Add(magicCodeTTL)); err != nil {
		return err
	}
	return u.mailer.SendMagicCode(ctx, mail.MagicCodeMessage{
		To:               email,
		Code:             code,
		Purpose:          string(p),
		ExpiresInMinutes: magicExpiresMin,
	})
}

// VerifyMagicCode verifies a magic code.
// Register creates the account if the email is new; login requires an existing account.
func (u *AuthUsecase) VerifyMagicCode(ctx context.Context, email, code string, purpose entity.MagicPurpose, meta SessionMeta) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	now := u.now()
	id, storedName, hash, attempts, err := u.codes.GetLatestActive(ctx, email, purpose, now)
	if err != nil {
		return nil, ErrCodeExpired
	}
	if attempts >= magicMaxAttempts {
		_ = u.codes.MarkConsumed(ctx, id, now)
		return nil, ErrTooManyAttempts
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(code)); err != nil {
		n, _ := u.codes.IncrementAttempts(ctx, id)
		if n >= magicMaxAttempts {
			_ = u.codes.MarkConsumed(ctx, id, now)
			return nil, ErrTooManyAttempts
		}
		return nil, ErrInvalidCredentials
	}
	if err := u.codes.MarkConsumed(ctx, id, now); err != nil {
		return nil, err
	}

	user, err := u.users.GetByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		// Login must not create accounts.
		if purpose == entity.PurposeLogin {
			return nil, ErrNoAccount
		}
		user = &entity.User{Email: email, Name: strings.TrimSpace(storedName)}
		if err := u.users.Create(ctx, user); err != nil {
			return nil, err
		}
		def := entity.DefaultSettings(user.ID)
		if err := u.settings.Upsert(ctx, &def); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if err := u.methods.Upsert(ctx, &entity.AuthMethod{
		UserID:         user.ID,
		Provider:       entity.ProviderEmail,
		ProviderUserID: email,
		Email:          email,
	}); err != nil {
		return nil, err
	}
	return u.issueSession(ctx, user, meta)
}

func (u *AuthUsecase) CompleteOAuthLogin(ctx context.Context, provider entity.Provider, providerUserID, email, name, avatarURL string, meta SessionMeta) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	var user *entity.User

	method, err := u.methods.GetByProvider(ctx, provider, providerUserID)
	if err == nil {
		user, err = u.users.GetByID(ctx, method.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		// Try match by verified provider email.
		existing, err := u.users.GetByEmail(ctx, email)
		if err == nil {
			user = existing
		} else {
			user = &entity.User{Email: email, Name: name, AvatarURL: avatarURL}
			if err := u.users.Create(ctx, user); err != nil {
				return nil, err
			}
			def := entity.DefaultSettings(user.ID)
			if err := u.settings.Upsert(ctx, &def); err != nil {
				return nil, err
			}
		}
	}

	if err := u.methods.Upsert(ctx, &entity.AuthMethod{
		UserID:         user.ID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		Email:          email,
	}); err != nil {
		return nil, err
	}
	if name != "" && user.Name == "" {
		user.Name = name
		_ = u.users.Update(ctx, user)
	}
	if avatarURL != "" && user.AvatarURL == "" {
		user.AvatarURL = avatarURL
		_ = u.users.Update(ctx, user)
	}
	return u.issueSession(ctx, user, meta)
}

func (u *AuthUsecase) LinkOAuth(ctx context.Context, userID uuid.UUID, provider entity.Provider, providerUserID, email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	existing, err := u.methods.GetByProvider(ctx, provider, providerUserID)
	if err == nil && existing.UserID != userID {
		return ErrConflict
	}
	if email != "" {
		owner, err := u.users.GetByEmail(ctx, email)
		if err == nil && owner.ID != userID {
			return ErrConflict
		}
	}
	return u.methods.Upsert(ctx, &entity.AuthMethod{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		Email:          email,
	})
}

func (u *AuthUsecase) UnlinkProvider(ctx context.Context, userID uuid.UUID, provider entity.Provider) error {
	n, err := u.methods.CountByUser(ctx, userID)
	if err != nil {
		return err
	}
	if n <= 1 {
		return ErrLastAuthMethod
	}
	return u.methods.Delete(ctx, userID, provider)
}

func (u *AuthUsecase) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return u.sessions.Revoke(ctx, sessionID, u.now())
}

func (u *AuthUsecase) SessionUser(ctx context.Context, token string) (*entity.User, *entity.Session, error) {
	hash := hashToken(token)
	sess, err := u.sessions.GetByTokenHash(ctx, hash)
	if err != nil {
		return nil, nil, ErrUnauthenticated
	}
	if !sess.Active() {
		return nil, nil, ErrUnauthenticated
	}
	user, err := u.users.GetByID(ctx, sess.UserID)
	if err != nil {
		return nil, nil, ErrUnauthenticated
	}
	// Throttle touch to once a minute.
	if time.Since(sess.LastSeenAt) > time.Minute {
		_ = u.sessions.Touch(ctx, sess.ID, u.now())
	}
	return user, sess, nil
}

func (u *AuthUsecase) issueSession(ctx context.Context, user *entity.User, meta SessionMeta) (*AuthResult, error) {
	raw, err := generateToken()
	if err != nil {
		return nil, err
	}
	newDevice := u.isNewDevice(ctx, user.ID, meta.UserAgent)
	location := u.geo.Resolve(ctx, meta.IP)
	sess := &entity.Session{
		TokenHash: hashToken(raw),
		UserID:    user.ID,
		UserAgent: meta.UserAgent,
		IP:        meta.IP,
		Location:  location,
	}
	if err := u.sessions.Create(ctx, sess); err != nil {
		return nil, err
	}
	_ = u.audit.Insert(ctx, &entity.LoginAudit{
		UserID:    user.ID,
		SessionID: sess.ID,
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
	})

	if newDevice && u.loginAlertsEnabled(ctx, user.ID) {
		_ = u.mailer.SendLoginAlert(ctx, mail.LoginAlertMessage{
			To:         user.Email,
			Device:     ua.Parse(meta.UserAgent).String(),
			IP:         meta.IP,
			Location:   location,
			LoggedInAt: u.now().Format(time.RFC3339),
		})
	}
	return &AuthResult{User: user, Session: sess, Token: raw}, nil
}

func (u *AuthUsecase) loginAlertsEnabled(ctx context.Context, userID uuid.UUID) bool {
	settings, err := u.settings.Get(ctx, userID)
	return err == nil && settings.LoginAlerts
}

func (u *AuthUsecase) isNewDevice(ctx context.Context, userID uuid.UUID, userAgent string) bool {
	prior, err := u.sessions.ListActive(ctx, userID)
	if err != nil {
		return true
	}
	current := ua.Parse(userAgent).String()
	for _, s := range prior {
		if ua.Parse(s.UserAgent).String() == current {
			return false
		}
	}
	return true
}

func (u *AuthUsecase) allowMagicRequest(key string) error {
	now := u.now()
	e := magicRate[key]
	if e == nil {
		e = &rateEntry{}
		magicRate[key] = e
	}
	cut := now.Add(-u.rateWindow)
	kept := e.times[:0]
	for _, t := range e.times {
		if t.After(cut) {
			kept = append(kept, t)
		}
	}
	e.times = kept
	if len(e.times) >= magicRatePerHour {
		return ErrRateLimited
	}
	e.times = append(e.times, now)
	return nil
}

func (u *AuthUsecase) allowResendCooldown(key string) bool {
	e := magicRate[key]
	if e == nil || len(e.times) == 0 {
		return true
	}
	last := e.times[len(e.times)-1]
	return u.now().Sub(last) >= magicResendWindow
}

func generateNumericCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
