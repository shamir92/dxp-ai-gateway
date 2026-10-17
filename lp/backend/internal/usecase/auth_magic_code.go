package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/mail"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
)

const (
	magicCodeTTL      = 10 * time.Minute
	magicMaxAttempts  = 5
	magicResendWindow = 60 * time.Second
	magicRatePerHour  = 5
	magicExpiresMin   = 10
)

// lastRequestAt is in-memory rate state for magic requests (per process).
type rateEntry struct{ times []time.Time }

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
	return u.mailer.SendMagicCode(ctx, mail.MagicCodeMessage{To: email, Code: code, Purpose: string(p), ExpiresInMinutes: magicExpiresMin})
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
	return u.mailer.SendMagicCode(ctx, mail.MagicCodeMessage{To: email, Code: code, Purpose: string(p), ExpiresInMinutes: magicExpiresMin})
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
	created := false
	if errors.Is(err, repository.ErrNotFound) {
		if purpose == entity.PurposeLogin {
			return nil, ErrNoAccount
		}
		user = &entity.User{Email: email, Name: strings.TrimSpace(storedName)}
		method := &entity.AuthMethod{UserID: user.ID, Provider: entity.ProviderEmail, ProviderUserID: email, Email: email}
		if err := u.createNewAccount(ctx, user, method); err != nil {
			return nil, err
		}
		created = true
	} else if err != nil {
		return nil, err
	}
	if !created {
		if err := u.methods.Upsert(ctx, &entity.AuthMethod{UserID: user.ID, Provider: entity.ProviderEmail, ProviderUserID: email, Email: email}); err != nil {
			return nil, err
		}
	}
	return u.issueSession(ctx, user, meta)
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
