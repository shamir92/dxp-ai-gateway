package usecase

import (
	"context"
	"errors"
	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
	"github.com/google/uuid"
)

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
