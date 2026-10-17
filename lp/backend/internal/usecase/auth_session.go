package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/mail"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/ua"
)

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
	sess := &entity.Session{TokenHash: hashToken(raw), UserID: user.ID, UserAgent: meta.UserAgent, IP: meta.IP, Location: location}
	if err := u.sessions.Create(ctx, sess); err != nil {
		return nil, err
	}
	_ = u.audit.Insert(ctx, &entity.LoginAudit{UserID: user.ID, SessionID: sess.ID, IP: meta.IP, UserAgent: meta.UserAgent})
	if newDevice && u.loginAlertsEnabled(ctx, user.ID) {
		_ = u.mailer.SendLoginAlert(ctx, mail.LoginAlertMessage{To: user.Email, Device: ua.Parse(meta.UserAgent).String(), IP: meta.IP, Location: location, LoggedInAt: u.now().Format(time.RFC3339)})
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
