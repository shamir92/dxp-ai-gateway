package entity

import (
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	ProviderEmail  Provider = "email"
	ProviderGoogle Provider = "google"
	ProviderGithub Provider = "github"
)

type AuthMethod struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Provider       Provider  `json:"provider"`
	ProviderUserID string    `json:"provider_user_id"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
}
