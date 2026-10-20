package credentials

import (
	"testing"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/worker/internal/entity"
)

func testKey() []byte { return []byte("0123456789abcdef0123456789abcdef") }

func TestNewCipherRejectsWrongKeySize(t *testing.T) {
	if _, err := NewCipher([]byte("too-short")); err == nil {
		t.Fatal("NewCipher accepted a key that is not 32 bytes")
	}
}

// The worker can only decrypt, so a round trip is impossible here. This vector
// pins the associated data and key version instead: it must match what the
// backend's Encrypt produces, or every stored password fails to open.
func TestAssociatedDataVector(t *testing.T) {
	userID := uuid.MustParse("11111111-2222-3333-4444-555555555555")

	for _, tc := range []struct {
		site entity.ExternalSite
		want string
	}{
		{entity.ExternalSiteTPAI, "v1:11111111-2222-3333-4444-555555555555:tpai"},
		{entity.ExternalSiteAPAI, "v1:11111111-2222-3333-4444-555555555555:apai"},
	} {
		if got := string(associatedData(userID, tc.site)); got != tc.want {
			t.Errorf("associatedData(%q) = %q, want %q", tc.site, got, tc.want)
		}
	}
}

func TestDecryptRejectsUnknownKeyVersion(t *testing.T) {
	cipher, err := NewCipher(testKey())
	if err != nil {
		t.Fatalf("NewCipher: %v", err)
	}
	nonce := make([]byte, cipher.aead.NonceSize())
	if _, err := cipher.Decrypt(uuid.New(), entity.ExternalSiteTPAI, []byte("x"), nonce, "v2"); err == nil {
		t.Fatal("Decrypt accepted an unsupported key version")
	}
}

func TestDecryptRejectsWrongNonceSize(t *testing.T) {
	cipher, err := NewCipher(testKey())
	if err != nil {
		t.Fatalf("NewCipher: %v", err)
	}
	if _, err := cipher.Decrypt(uuid.New(), entity.ExternalSiteTPAI, []byte("x"), []byte{1, 2}, KeyVersion); err == nil {
		t.Fatal("Decrypt accepted a malformed nonce")
	}
}

func TestDecryptFailsForAnotherSite(t *testing.T) {
	// A ciphertext sealed for one site must not open for a different user or
	// site, which is what the associated data buys us.
	cipher, err := NewCipher(testKey())
	if err != nil {
		t.Fatalf("NewCipher: %v", err)
	}
	nonce := make([]byte, cipher.aead.NonceSize())
	if _, err := cipher.Decrypt(uuid.New(), entity.ExternalSiteAPAI, []byte("not-real"), nonce, KeyVersion); err == nil {
		t.Fatal("Decrypt opened a ciphertext that was never sealed")
	}
}
