package credentials

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
)

const generatedPasswordLength = 32

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
	displayName := strings.TrimSpace(user.Name)
	if displayName == "" {
		displayName = user.Email
	}
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
			Username:             strings.TrimSpace(user.Email),
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
