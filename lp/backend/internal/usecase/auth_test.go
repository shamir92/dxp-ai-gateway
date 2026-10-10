package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/mail"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
	"github.com/dxpglobal/ai-gateway/backend/internal/usecase"
)

type memUsers struct{ byEmail map[string]*entity.User; byID map[uuid.UUID]*entity.User }

func newMemUsers() *memUsers {
	return &memUsers{byEmail: map[string]*entity.User{}, byID: map[uuid.UUID]*entity.User{}}
}
func (m *memUsers) Create(_ context.Context, u *entity.User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	m.byEmail[u.Email] = u
	m.byID[u.ID] = u
	return nil
}
func (m *memUsers) GetByID(_ context.Context, id uuid.UUID) (*entity.User, error) {
	if u, ok := m.byID[id]; ok {
		return u, nil
	}
	return nil, repository.ErrNotFound
}
func (m *memUsers) GetByEmail(_ context.Context, email string) (*entity.User, error) {
	if u, ok := m.byEmail[email]; ok {
		return u, nil
	}
	return nil, repository.ErrNotFound
}
func (m *memUsers) Update(_ context.Context, u *entity.User) error {
	m.byEmail[u.Email] = u
	m.byID[u.ID] = u
	return nil
}

type memMethods struct{ items []entity.AuthMethod }

func (m *memMethods) Upsert(_ context.Context, a *entity.AuthMethod) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	for i, x := range m.items {
		if x.Provider == a.Provider && x.ProviderUserID == a.ProviderUserID {
			m.items[i] = *a
			return nil
		}
	}
	m.items = append(m.items, *a)
	return nil
}
func (m *memMethods) ListByUser(_ context.Context, userID uuid.UUID) ([]entity.AuthMethod, error) {
	var out []entity.AuthMethod
	for _, x := range m.items {
		if x.UserID == userID {
			out = append(out, x)
		}
	}
	return out, nil
}
func (m *memMethods) GetByProvider(_ context.Context, p entity.Provider, id string) (*entity.AuthMethod, error) {
	for i := range m.items {
		if m.items[i].Provider == p && m.items[i].ProviderUserID == id {
			return &m.items[i], nil
		}
	}
	return nil, repository.ErrNotFound
}
func (m *memMethods) Delete(_ context.Context, userID uuid.UUID, p entity.Provider) error {
	out := m.items[:0]
	for _, x := range m.items {
		if !(x.UserID == userID && x.Provider == p) {
			out = append(out, x)
		}
	}
	m.items = out
	return nil
}
func (m *memMethods) CountByUser(_ context.Context, userID uuid.UUID) (int, error) {
	n := 0
	for _, x := range m.items {
		if x.UserID == userID {
			n++
		}
	}
	return n, nil
}

type memSessions struct{ byHash map[string]*entity.Session }

func newMemSessions() *memSessions { return &memSessions{byHash: map[string]*entity.Session{}} }
func (m *memSessions) Create(_ context.Context, s *entity.Session) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	cp := *s
	m.byHash[s.TokenHash] = &cp
	return nil
}
func (m *memSessions) GetByTokenHash(_ context.Context, hash string) (*entity.Session, error) {
	if s, ok := m.byHash[hash]; ok {
		cp := *s
		return &cp, nil
	}
	return nil, repository.ErrNotFound
}
func (m *memSessions) Touch(_ context.Context, id uuid.UUID, at time.Time) error {
	for _, s := range m.byHash {
		if s.ID == id {
			s.LastSeenAt = at
		}
	}
	return nil
}
func (m *memSessions) Revoke(_ context.Context, id uuid.UUID, at time.Time) error {
	for _, s := range m.byHash {
		if s.ID == id {
			t := at
			s.RevokedAt = &t
		}
	}
	return nil
}
func (m *memSessions) RevokeAllForUser(_ context.Context, userID uuid.UUID, at time.Time) error {
	for _, s := range m.byHash {
		if s.UserID == userID {
			t := at
			s.RevokedAt = &t
		}
	}
	return nil
}
func (m *memSessions) ListActive(_ context.Context, userID uuid.UUID) ([]entity.Session, error) {
	var out []entity.Session
	for _, s := range m.byHash {
		if s.UserID == userID && s.RevokedAt == nil {
			out = append(out, *s)
		}
	}
	return out, nil
}

type memSettings struct{ byID map[uuid.UUID]*entity.UserSettings }

func newMemSettings() *memSettings { return &memSettings{byID: map[uuid.UUID]*entity.UserSettings{}} }
func (m *memSettings) Get(_ context.Context, id uuid.UUID) (*entity.UserSettings, error) {
	if s, ok := m.byID[id]; ok {
		cp := *s
		return &cp, nil
	}
	return nil, repository.ErrNotFound
}
func (m *memSettings) Upsert(_ context.Context, s *entity.UserSettings) error {
	cp := *s
	m.byID[s.UserID] = &cp
	return nil
}

type memCodes struct {
	last   *entity.MagicPurpose
	store  map[uuid.UUID]codeRow
	seq    int
}
type codeRow struct {
	id        uuid.UUID
	email     string
	name      string
	hash      string
	purpose   entity.MagicPurpose
	expiresAt time.Time
	consumed  bool
	attempts  int
	createdAt time.Time
}

func newMemCodes() *memCodes { return &memCodes{store: map[uuid.UUID]codeRow{}} }
func (m *memCodes) Insert(_ context.Context, email, name, codeHash string, purpose entity.MagicPurpose, expiresAt time.Time) error {
	m.seq++
	id := uuid.New()
	m.store[id] = codeRow{id: id, email: email, name: name, hash: codeHash, purpose: purpose, expiresAt: expiresAt, createdAt: time.Now()}
	return nil
}
func (m *memCodes) GetLatestActive(_ context.Context, email string, purpose entity.MagicPurpose, now time.Time) (uuid.UUID, string, string, int, error) {
	var best *codeRow
	for _, r := range m.store {
		r := r
		if r.email == email && r.purpose == purpose && !r.consumed && r.expiresAt.After(now) {
			if best == nil || r.createdAt.After(best.createdAt) {
				best = &r
			}
		}
	}
	if best == nil {
		return uuid.Nil, "", "", 0, repository.ErrNotFound
	}
	return best.id, best.name, best.hash, best.attempts, nil
}
func (m *memCodes) InvalidateForEmail(_ context.Context, email string, purpose entity.MagicPurpose) error {
	for id, r := range m.store {
		if r.email == email && r.purpose == purpose && !r.consumed {
			r.consumed = true
			m.store[id] = r
		}
	}
	return nil
}
func (m *memCodes) IncrementAttempts(_ context.Context, id uuid.UUID) (int, error) {
	r := m.store[id]
	r.attempts++
	m.store[id] = r
	return r.attempts, nil
}
func (m *memCodes) MarkConsumed(_ context.Context, id uuid.UUID, _ time.Time) error {
	r := m.store[id]
	r.consumed = true
	m.store[id] = r
	return nil
}

type memAudit struct{}

func (m *memAudit) Insert(_ context.Context, _ *entity.LoginAudit) error { return nil }
func (m *memAudit) ListByUser(_ context.Context, _ uuid.UUID, _ int) ([]entity.LoginAudit, error) {
	return nil, nil
}

type captureMailer struct{ last *mail.MagicCodeMessage }

func (m *captureMailer) SendMagicCode(_ context.Context, msg mail.MagicCodeMessage) error {
	m.last = &msg
	return nil
}

func (m *captureMailer) SendLoginAlert(_ context.Context, _ mail.LoginAlertMessage) error {
	return nil
}

func TestMagicRegisterAndLogin(t *testing.T) {
	users := newMemUsers()
	methods := &memMethods{}
	sessions := newMemSessions()
	settings := newMemSettings()
	codes := newMemCodes()
	audit := &memAudit{}
	mailer := &captureMailer{}
	auth := usecase.NewAuthUsecase(users, methods, sessions, settings, codes, audit, mailer, nil)

	ctx := context.Background()
	if err := auth.RequestMagicCode(ctx, "Ada@example.com", "register", "Ada Lovelace"); err != nil {
		t.Fatalf("request: %v", err)
	}
	if mailer.last == nil || mailer.last.Code == "" {
		t.Fatal("expected mail with code")
	}

	// wrong code
	if _, err := auth.VerifyMagicCode(ctx, "ada@example.com", "000000", entity.PurposeRegister, usecase.SessionMeta{}); err == nil {
		t.Fatal("expected invalid code")
	}

	res, err := auth.VerifyMagicCode(ctx, "ada@example.com", mailer.last.Code, entity.PurposeRegister, usecase.SessionMeta{})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if res.User.Name != "Ada Lovelace" {
		t.Fatalf("name = %q", res.User.Name)
	}
	if res.Token == "" {
		t.Fatal("expected session token")
	}

	user, sess, err := auth.SessionUser(ctx, res.Token)
	if err != nil || user.ID != res.User.ID || sess.ID != res.Session.ID {
		t.Fatalf("session user: err=%v user=%v sess=%v", err, user, sess)
	}

	// login again as existing user — name not overwritten
	if err := auth.RequestMagicCode(ctx, "ada@example.com", "login", "Other"); err != nil {
		t.Fatalf("login request: %v", err)
	}
	res2, err := auth.VerifyMagicCode(ctx, "ada@example.com", mailer.last.Code, entity.PurposeLogin, usecase.SessionMeta{})
	if err != nil {
		t.Fatalf("login verify: %v", err)
	}
	if res2.User.Name != "Ada Lovelace" {
		t.Fatalf("name overwritten: %q", res2.User.Name)
	}
}

func TestLoginDoesNotCreateAccount(t *testing.T) {
	users := newMemUsers()
	mailer := &captureMailer{}
	auth := usecase.NewAuthUsecase(users, &memMethods{}, newMemSessions(), newMemSettings(), newMemCodes(), &memAudit{}, mailer, nil)
	ctx := context.Background()

	err := auth.RequestMagicCode(ctx, "nobody@example.com", "login", "")
	if !errors.Is(err, usecase.ErrNoAccount) {
		t.Fatalf("login request for unknown email: %v", err)
	}
	if mailer.last != nil {
		t.Fatal("should not email a code for unknown login email")
	}

	// Register creates the account, then login works.
	if err := auth.RequestMagicCode(ctx, "nobody@example.com", "register", "Nobody"); err != nil {
		t.Fatalf("register request: %v", err)
	}
	if mailer.last == nil {
		t.Fatal("register should email a code")
	}
	// Duplicate register rejected.
	if err := auth.RequestMagicCode(ctx, "nobody@example.com", "register", "Nobody"); !errors.Is(err, usecase.ErrAccountExists) {
		// user not created until verify — so this still sends a code path.
		// Verify first, then expect account_exists.
	}
	res, err := auth.VerifyMagicCode(ctx, "nobody@example.com", mailer.last.Code, entity.PurposeRegister, usecase.SessionMeta{})
	if err != nil {
		t.Fatalf("register verify: %v", err)
	}
	if res.User == nil {
		t.Fatal("expected user")
	}
	if err := auth.RequestMagicCode(ctx, "nobody@example.com", "register", "Nobody"); !errors.Is(err, usecase.ErrAccountExists) {
		t.Fatalf("duplicate register: %v", err)
	}
	if err := auth.RequestMagicCode(ctx, "nobody@example.com", "login", ""); err != nil {
		t.Fatalf("login request after register: %v", err)
	}
}

func TestUnlinkLastMethod(t *testing.T) {
	users := newMemUsers()
	methods := &memMethods{}
	auth := usecase.NewAuthUsecase(users, methods, newMemSessions(), newMemSettings(), newMemCodes(), &memAudit{}, &captureMailer{}, nil)
	ctx := context.Background()
	u := &entity.User{Email: "a@b.com"}
	_ = users.Create(ctx, u)
	_ = methods.Upsert(ctx, &entity.AuthMethod{UserID: u.ID, Provider: entity.ProviderEmail, ProviderUserID: "a@b.com", Email: "a@b.com"})
	if err := auth.UnlinkProvider(ctx, u.ID, entity.ProviderEmail); err != usecase.ErrLastAuthMethod {
		t.Fatalf("err = %v", err)
	}
	_ = methods.Upsert(ctx, &entity.AuthMethod{UserID: u.ID, Provider: entity.ProviderGoogle, ProviderUserID: "g1", Email: "a@b.com"})
	if err := auth.UnlinkProvider(ctx, u.ID, entity.ProviderEmail); err != nil {
		t.Fatalf("unlink: %v", err)
	}
}
