package credentials

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
)

const generatedPasswordLength = 32

// Upstream new-api validates both `username` and `display_name` with `max=20`
// (model.UserNameMaxLength), so a gateway email longer than that cannot be used
// verbatim as the remote username — a 22-char address is rejected outright.
const (
	remoteNameMaxLength = 20
	remoteHashChars     = 6
)

// RemoteUsername derives a stable, upstream-safe username from a gateway email.
// It keeps a readable prefix and appends a short hash of the full address so two
// long local parts cannot collide.
func RemoteUsername(email string) string {
	normalized := strings.ToLower(strings.TrimSpace(email))
	local, _, _ := strings.Cut(normalized, "@")

	prefix := make([]byte, 0, len(local))
	for i := 0; i < len(local); i++ {
		switch c := local[i]; {
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9', c == '.', c == '_', c == '-':
			prefix = append(prefix, c)
		}
	}
	if maxPrefix := remoteNameMaxLength - remoteHashChars - 1; len(prefix) > maxPrefix {
		prefix = prefix[:maxPrefix]
	}

	sum := sha256.Sum256([]byte(normalized))
	return string(prefix) + "-" + hex.EncodeToString(sum[:])[:remoteHashChars]
}

// remoteDisplayName prefers the user's name and falls back to the email's local
// part (which reads better than a truncated address), capped to the same limit.
func remoteDisplayName(user *entity.User) string {
	name := strings.TrimSpace(user.Name)
	if name == "" {
		local, _, _ := strings.Cut(strings.ToLower(strings.TrimSpace(user.Email)), "@")
		name = local
	}
	runes := []rune(name)
	if len(runes) > remoteNameMaxLength {
		return string(runes[:remoteNameMaxLength])
	}
	return name
}

type AccountFactory struct {
	cipher *Cipher
	sites  []entity.ExternalSite
}

func NewAccountFactory(cipher *Cipher, sites []entity.ExternalSite) (*AccountFactory, error) {
	if cipher == nil {
		return nil, fmt.Errorf("credential cipher is required")
	}
	if len(sites) != 2 || sites[0] == sites[1] {
		return nil, fmt.Errorf("exactly two distinct external sites are required")
	}
	for _, site := range sites {
		if site != entity.ExternalSiteTPAI && site != entity.ExternalSiteAPAI {
			return nil, fmt.Errorf("unsupported external site")
		}
	}
	return &AccountFactory{cipher: cipher, sites: append([]entity.ExternalSite(nil), sites...)}, nil
}

func (f *AccountFactory) Build(user *entity.User) ([]entity.ExternalSiteAccount, error) {
	if user == nil || user.ID == uuid.Nil || strings.TrimSpace(user.Email) == "" {
		return nil, fmt.Errorf("user ID and email are required for external accounts")
	}
	displayName := remoteDisplayName(user)
	username := RemoteUsername(user.Email)
	accounts := make([]entity.ExternalSiteAccount, 0, len(f.sites))
	for _, site := range f.sites {
		password, err := randomPassword(generatedPasswordLength)
		if err != nil {
			return nil, err
		}
		ciphertext, nonce, version, err := f.cipher.Encrypt(user.ID, site, password)
		clear(password)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, entity.ExternalSiteAccount{
			ID:                   uuid.New(),
			UserID:               user.ID,
			Site:                 site,
			Username:             username,
			DisplayName:          displayName,
			PasswordCiphertext:   ciphertext,
			PasswordNonce:        nonce,
			EncryptionKeyVersion: version,
			Status:               entity.ExternalAccountPending,
			NextAttemptAt:        time.Now().UTC(),
		})
	}
	return accounts, nil
}

func randomPassword(length int) ([]byte, error) {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	password := make([]byte, length)
	max := big.NewInt(int64(len(alphabet)))
	for i := range password {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			clear(password)
			return nil, fmt.Errorf("generate external account password: %w", err)
		}
		password[i] = alphabet[n.Int64()]
	}
	return password, nil
}
