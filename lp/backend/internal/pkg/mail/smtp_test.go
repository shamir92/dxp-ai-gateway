package mail

import "testing"

func TestAddressStripsDisplayName(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"bare", "noreply@example.com", "noreply@example.com"},
		{"display name", "DXP AI Gateway <noreply@example.com>", "noreply@example.com"},
		{"quoted display name", "\"DXP, Inc.\" <noreply@example.com>", "noreply@example.com"},
		{"leading space", " noreply@example.com", "noreply@example.com"},
		{"space inside brackets", "DXP < noreply@example.com >", "noreply@example.com"},
		{"unterminated", "DXP <noreply@example.com", "noreply@example.com"},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := address(tc.in); got != tc.want {
				t.Fatalf("address(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestEnvelopeFromIsBareAddress(t *testing.T) {
	m := &SMTPMailer{cfg: SMTPConfig{From: "noreply@example.com", FromName: "DXP AI Gateway"}}

	// The header keeps the display name, the envelope must not.
	from := m.headerFrom()
	if from != "DXP AI Gateway <noreply@example.com>" {
		t.Fatalf("header From = %q", from)
	}
	if got := address(m.cfg.From); got != "noreply@example.com" {
		t.Fatalf("envelope From = %q, want bare address", got)
	}
}
