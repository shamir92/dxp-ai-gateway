package credentials

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
)

const KeyVersion = "v1"

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

func (c *Cipher) Encrypt(userID uuid.UUID, site entity.ExternalSite, plaintext []byte) (ciphertext, nonce []byte, keyVersion string, err error) {
	nonce = make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, "", fmt.Errorf("generate encryption nonce: %w", err)
	}
	ciphertext = c.aead.Seal(nil, nonce, plaintext, associatedData(userID, site))
	return ciphertext, nonce, KeyVersion, nil
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

func associatedData(userID uuid.UUID, site entity.ExternalSite) []byte {
	return []byte(KeyVersion + ":" + userID.String() + ":" + string(site))
}
