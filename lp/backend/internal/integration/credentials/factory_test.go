package credentials

import (
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
)

func TestRemoteUsernameFitsUpstreamLimit(t *testing.T) {
	// new-api validates username and display_name with max=20, so every derived
	// handle must fit, whatever the email looks like.
	emails := []string{
		"shamirhusein@gmail.com", // 22 chars: rejected upstream before this fix
		"a-very-long-local-part-indeed@example.com",
		"shamir+3@gmail.com",
		"SHAMIRHUSEIN@GMAIL.COM",
		"x@y.io",
		"weird!#$local@example.com",
	}
	for _, email := range emails {
		got := RemoteUsername(email)
		if len(got) > remoteNameMaxLength {
			t.Errorf("RemoteUsername(%q) = %q, %d chars, want <= %d", email, got, len(got), remoteNameMaxLength)
		}
		if got == "" {
			t.Errorf("RemoteUsername(%q) is empty", email)
		}
	}
}

func TestRemoteUsernameIsDeterministicAndCaseInsensitive(t *testing.T) {
	base := RemoteUsername("shamirhusein@gmail.com")
	if again := RemoteUsername("shamirhusein@gmail.com"); again != base {
		t.Errorf("not deterministic: %q vs %q", base, again)
	}
	if upper := RemoteUsername("SHAMIRHUSEIN@GMAIL.COM"); upper != base {
		t.Errorf("case changed the handle: %q vs %q", upper, base)
	}
	if padded := RemoteUsername("  shamirhusein@gmail.com  "); padded != base {
		t.Errorf("whitespace changed the handle: %q vs %q", padded, base)
	}
}

func TestRemoteUsernameKeepsLongLocalPartsApart(t *testing.T) {
	// Both truncate to the same readable prefix, so the hash must disambiguate.
	a := RemoteUsername("shamirhuseinone@gmail.com")
	b := RemoteUsername("shamirhuseintwo@gmail.com")
	if a == b {
		t.Errorf("collision: both %q", a)
	}
}

func TestRemoteDisplayNameIsCapped(t *testing.T) {
	long := &entity.User{Email: "shamirhusein@gmail.com", Name: "Shamir Husein FromDxp Indonesia"}
	if got := remoteDisplayName(long); len([]rune(got)) > remoteNameMaxLength {
		t.Errorf("display name %q exceeds %d", got, remoteNameMaxLength)
	}
	// With no name, prefer the local part over the full address.
	unnamed := &entity.User{Email: "shamirhusein@gmail.com"}
	if got := remoteDisplayName(unnamed); got != "shamirhusein" {
		t.Errorf("display name = %q, want the local part", got)
	}
}

func TestBuildProducesUpstreamSafeIdentities(t *testing.T) {
	cipher, err := NewCipher([]byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatalf("cipher: %v", err)
	}
	factory, err := NewAccountFactory(cipher, []entity.ExternalSite{entity.ExternalSiteTPAI, entity.ExternalSiteAPAI})
	if err != nil {
		t.Fatalf("factory: %v", err)
	}

	accounts, err := factory.Build(&entity.User{ID: uuid.New(), Email: "shamirhusein@gmail.com"})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("accounts = %d, want 2 (one per site)", len(accounts))
	}
	for _, account := range accounts {
		if len(account.Username) > remoteNameMaxLength {
			t.Errorf("username %q exceeds %d", account.Username, remoteNameMaxLength)
		}
		if len([]rune(account.DisplayName)) > remoteNameMaxLength {
			t.Errorf("display name %q exceeds %d", account.DisplayName, remoteNameMaxLength)
		}
		if strings.Contains(account.Username, "@") {
			t.Errorf("username %q still looks like an email", account.Username)
		}
		if len(account.PasswordCiphertext) == 0 {
			t.Error("password was not encrypted")
		}
	}
}
