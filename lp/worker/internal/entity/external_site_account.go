package entity

import (
	"time"

	"github.com/google/uuid"
)

type ExternalSite string

const (
	ExternalSiteTPAI ExternalSite = "tpai"
	ExternalSiteAPAI ExternalSite = "apai"
)

type ExternalSiteAccountStatus string

const (
	ExternalAccountPending     ExternalSiteAccountStatus = "pending"
	ExternalAccountRetrying    ExternalSiteAccountStatus = "retrying"
	ExternalAccountProvisioned ExternalSiteAccountStatus = "provisioned"
	ExternalAccountFailed      ExternalSiteAccountStatus = "failed"
)

type ExternalSiteAccount struct {
	ID                   uuid.UUID
	UserID               uuid.UUID
	Site                 ExternalSite
	Username             string
	DisplayName          string
	PasswordCiphertext   []byte
	PasswordNonce        []byte
	EncryptionKeyVersion string
	RemoteUserID         *string
	Status               ExternalSiteAccountStatus
	AttemptCount         int
	NextAttemptAt        time.Time
	LeaseToken           uuid.UUID
	LeaseExpiresAt       *time.Time
	LastError            string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ProvisionedAt        *time.Time
}
