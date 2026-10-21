package credentials

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/worker/internal/entity"
)

const KeyVersion = "v1"

// Cipher decrypts the per-user site password the backend sealed at signup.
//
// The algorithm, the key version and associatedData below must stay identical to
// the backend's cipher: every stored ciphertext was sealed with exactly this key
// derivation, so any change silently breaks all existing rows.
type Cipher struct {
	aead cipher.AEAD
}

func NewCipher(key []byte) (*Cipher, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("AES-256-GCM key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create AES-GCM cipher: %w", err)
	}
	return &Cipher{aead: aead}, nil
}

func (c *Cipher) Decrypt(userID uuid.UUID, site entity.ExternalSite, ciphertext, nonce []byte, keyVersion string) ([]byte, error) {
	if keyVersion != KeyVersion {
		return nil, fmt.Errorf("unsupported encryption key version")
	}
	if len(nonce) != c.aead.NonceSize() {
		return nil, fmt.Errorf("invalid encrypted credential")
	}
	plaintext, err := c.aead.Open(nil, nonce, ciphertext, associatedData(userID, site))
	if err != nil {
		return nil, fmt.Errorf("decrypt external account password")
	}
	return plaintext, nil
}

// associatedData is the password binding. Do not change it: every stored password
// was sealed with exactly this value.
func associatedData(userID uuid.UUID, site entity.ExternalSite) []byte {
	return []byte(KeyVersion + ":" + userID.String() + ":" + string(site))
}
