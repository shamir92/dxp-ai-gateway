package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
)

func (u *AuthUsecase) CompleteOAuthLogin(ctx context.Context, provider entity.Provider, providerUserID, email, name, avatarURL string, meta SessionMeta) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	var user *entity.User
	method, err := u.methods.GetByProvider(ctx, provider, providerUserID)
	created := false
	if err == nil {
		user, err = u.users.GetByID(ctx, method.UserID)
		if err != nil {
			return nil, err
		}
	} else if errors.Is(err, repository.ErrNotFound) {
		existing, err := u.users.GetByEmail(ctx, email)
		if err == nil {
			user = existing
		} else if errors.Is(err, repository.ErrNotFound) {
			if email == "" {
				return nil, fmt.Errorf("%w: email required for account creation", ErrValidation)
			}
			user = &entity.User{Email: email, Name: name, AvatarURL: avatarURL}
			newMethod := &entity.AuthMethod{UserID: user.ID, Provider: provider, ProviderUserID: providerUserID, Email: email}
			if err := u.createNewAccount(ctx, user, newMethod); err != nil {
				return nil, err
			}
			created = true
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
	if !created {
		if err := u.methods.Upsert(ctx, &entity.AuthMethod{UserID: user.ID, Provider: provider, ProviderUserID: providerUserID, Email: email}); err != nil {
			return nil, err
		}
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
