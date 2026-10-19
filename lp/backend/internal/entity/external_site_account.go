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

// AccessTokenStatus tracks whether the stored personal access token is usable.
type AccessTokenStatus string

const (
	AccessTokenNone    AccessTokenStatus = "none"
	AccessTokenActive  AccessTokenStatus = "active"
	AccessTokenInvalid AccessTokenStatus = "invalid"
)

// StoredAccessToken is the encrypted personal access token ("nap_…") for one
// (user, site). The plaintext never leaves the backend process.
type StoredAccessToken struct {
	UserID     uuid.UUID
	Site       ExternalSite
	Ciphertext []byte
	Nonce      []byte
	KeyVersion string
	ExpiresAt  *time.Time // nil means the token never expires.
	RemoteID   string
	Status     AccessTokenStatus
}
