package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
)

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
	if err := u.methods.Upsert(ctx, &entity.AuthMethod{UserID: user.ID, Provider: provider, ProviderUserID: providerUserID, Email: email}); err != nil {
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
	return u.methods.Upsert(ctx, &entity.AuthMethod{UserID: userID, Provider: provider, ProviderUserID: providerUserID, Email: email})
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
