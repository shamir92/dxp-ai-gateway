package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/ua"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
)

type ProfileUsecase struct {
	users    repository.UserRepo
	methods  repository.AuthMethodRepo
	sessions repository.SessionRepo
	settings repository.SettingsRepo
	now      func() time.Time
}

func NewProfileUsecase(
	users repository.UserRepo,
	methods repository.AuthMethodRepo,
	sessions repository.SessionRepo,
	settings repository.SettingsRepo,
) *ProfileUsecase {
	return &ProfileUsecase{
		users:    users,
		methods:  methods,
		sessions: sessions,
		settings: settings,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

type ProfileView struct {
	User         *entity.User         `json:"user"`
	Connections  []ConnectionView     `json:"connections"`
	SessionCount int                  `json:"session_count"`
	Settings     *entity.UserSettings `json:"settings"`
}

type SessionView struct {
	ID         uuid.UUID `json:"id"`
	Device     string    `json:"device"`
	UserAgent  string    `json:"user_agent"`
	IP         string    `json:"ip"`
	Location   string    `json:"location"`
	CreatedAt  time.Time `json:"created_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
	Current    bool      `json:"current"`
}

type ConnectionView struct {
	Provider  entity.Provider `json:"provider"`
	Connected bool            `json:"connected"`
	Email     string          `json:"email"`
}

func (u *ProfileUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*ProfileView, error) {
	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	methods, err := u.methods.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	sessions, err := u.sessions.ListActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	settings, err := u.settings.Get(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		def := entity.DefaultSettings(userID)
		settings = &def
		_ = u.settings.Upsert(ctx, settings)
	} else if err != nil {
		return nil, err
	}
	return &ProfileView{
		User:         user,
		Connections:  connectionViews(methods, user.Email),
		SessionCount: len(sessions),
		Settings:     settings,
	}, nil
}

func connectionViews(methods []entity.AuthMethod, fallbackEmail string) []ConnectionView {
	byProv := map[entity.Provider]string{}
	for _, m := range methods {
		byProv[m.Provider] = m.Email
		if byProv[m.Provider] == "" {
			byProv[m.Provider] = fallbackEmail
		}
	}
	// Always show the three login methods. Email is connected when the user
	// signed up / signed in with a magic code (auth_methods provider=email).
	providers := []entity.Provider{entity.ProviderEmail, entity.ProviderGoogle, entity.ProviderGithub}
	out := make([]ConnectionView, 0, len(providers))
	for _, p := range providers {
		email, ok := byProv[p]
		if !ok && p == entity.ProviderEmail && fallbackEmail != "" {
			// Primary email is the passwordless identity even if the row is missing.
			out = append(out, ConnectionView{Provider: p, Connected: true, Email: fallbackEmail})
			continue
		}
		out = append(out, ConnectionView{Provider: p, Connected: ok, Email: email})
	}
	return out
}

func (u *ProfileUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email *string) (*entity.User, error) {
	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if name != nil {
		user.Name = strings.TrimSpace(*name)
	}
	if email != nil {
		e := strings.TrimSpace(strings.ToLower(*email))
		if e == "" || !strings.Contains(e, "@") {
			return nil, ErrValidation
		}
		existing, err := u.users.GetByEmail(ctx, e)
		if err == nil && existing.ID != userID {
			return nil, ErrConflict
		}
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
		user.Email = e
	}
	if err := u.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *ProfileUsecase) ListSessions(ctx context.Context, userID, currentID uuid.UUID) ([]SessionView, error) {
	list, err := u.sessions.ListActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]SessionView, 0, len(list))
	for _, s := range list {
		out = append(out, SessionView{
			ID:         s.ID,
			Device:     ua.Parse(s.UserAgent).String(),
			UserAgent:  s.UserAgent,
			IP:         s.IP,
			Location:   s.Location,
			CreatedAt:  s.CreatedAt,
			LastSeenAt: s.LastSeenAt,
			Current:    s.ID == currentID,
		})
	}
	return out, nil
}

func (u *ProfileUsecase) RevokeSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	list, err := u.sessions.ListActive(ctx, userID)
	if err != nil {
		return err
	}
	for _, s := range list {
		if s.ID == sessionID {
			return u.sessions.Revoke(ctx, sessionID, u.now())
		}
	}
	return ErrNotFound
}

func (u *ProfileUsecase) RevokeAllSessions(ctx context.Context, userID uuid.UUID) error {
	return u.sessions.RevokeAllForUser(ctx, userID, u.now())
}

func (u *ProfileUsecase) ListConnections(ctx context.Context, userID uuid.UUID) ([]ConnectionView, error) {
	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	methods, err := u.methods.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return connectionViews(methods, user.Email), nil
}

func (u *ProfileUsecase) UpdateSecurity(ctx context.Context, userID uuid.UUID, loginAlerts bool) (*entity.UserSettings, error) {
	s, err := u.settings.Get(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		def := entity.DefaultSettings(userID)
		s = &def
	} else if err != nil {
		return nil, err
	}
	s.LoginAlerts = loginAlerts
	if err := u.settings.Upsert(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}

type SettingsPatch struct {
	NotifyBalance  *bool `json:"notify_balance"`
	NotifyUsage    *bool `json:"notify_usage"`
	NotifySecurity *bool `json:"notify_security"`
	NotifyProduct  *bool `json:"notify_product"`
	LoginAlerts    *bool `json:"login_alerts"`
}

type SettingsUsecase struct {
	settings repository.SettingsRepo
	users    repository.UserRepo
}

func NewSettingsUsecase(settings repository.SettingsRepo, users repository.UserRepo) *SettingsUsecase {
	return &SettingsUsecase{settings: settings, users: users}
}

func (u *SettingsUsecase) Get(ctx context.Context, userID uuid.UUID) (*entity.UserSettings, error) {
	s, err := u.settings.Get(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		def := entity.DefaultSettings(userID)
		if err := u.settings.Upsert(ctx, &def); err != nil {
			return nil, err
		}
		return &def, nil
	}
	return s, err
}

func (u *SettingsUsecase) Patch(ctx context.Context, userID uuid.UUID, p SettingsPatch) (*entity.UserSettings, error) {
	s, err := u.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p.NotifyBalance != nil {
		s.NotifyBalance = *p.NotifyBalance
	}
	if p.NotifyUsage != nil {
		s.NotifyUsage = *p.NotifyUsage
	}
	if p.NotifySecurity != nil {
		s.NotifySecurity = *p.NotifySecurity
	}
	if p.NotifyProduct != nil {
		s.NotifyProduct = *p.NotifyProduct
	}
	if p.LoginAlerts != nil {
		s.LoginAlerts = *p.LoginAlerts
	}
	if err := u.settings.Upsert(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}
