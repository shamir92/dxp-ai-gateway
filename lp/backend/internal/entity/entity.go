package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) Initials() string {
	name := u.Name
	if name == "" {
		name = u.Email
	}
	out := []byte{}
	for _, part := range splitWords(name) {
		if len(part) > 0 {
			out = append(out, upper(part[0]))
		}
		if len(out) >= 2 {
			break
		}
	}
	if len(out) == 0 {
		return "?"
	}
	return string(out)
}

func splitWords(s string) []string {
	var words []string
	cur := []byte{}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '.' || c == '_' || c == '-' || c == '@' {
			if len(cur) > 0 {
				words = append(words, string(cur))
				cur = cur[:0]
			}
			continue
		}
		cur = append(cur, c)
	}
	if len(cur) > 0 {
		words = append(words, string(cur))
	}
	return words
}

func upper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 32
	}
	return c
}

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

type Session struct {
	ID         uuid.UUID  `json:"id"`
	TokenHash  string     `json:"-"`
	UserID     uuid.UUID  `json:"user_id"`
	UserAgent  string     `json:"user_agent"`
	IP         string     `json:"ip"`
	Location   string     `json:"location"`
	CreatedAt  time.Time  `json:"created_at"`
	LastSeenAt time.Time  `json:"last_seen_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

func (s Session) Active() bool {
	return s.RevokedAt == nil
}

type UserSettings struct {
	UserID           uuid.UUID `json:"user_id"`
	DefaultRegion    string    `json:"default_region"`
	DefaultModel     string    `json:"default_model"`
	NotifyBalance    bool      `json:"notify_balance"`
	NotifyUsage      bool      `json:"notify_usage"`
	NotifySecurity   bool      `json:"notify_security"`
	NotifyProduct    bool      `json:"notify_product"`
	StreamDefault    bool      `json:"stream_default"`
	StoreRequestLogs bool      `json:"store_request_logs"`
	LoginAlerts      bool      `json:"login_alerts"`
	UpdatedAt        time.Time `json:"updated_at"`
}

const (
	RegionSGP = "sgp"
	RegionUS  = "us"
)

var AllowedModels = map[string]bool{
	"deepseek-v4.1-flash": true,
	"deepseek-v4-flash":   true,
	"mimo-v2.6-pro":       true,
	"mimo-v2.6-flash":     true,
	"glm-5.3-flash":       true,
	"glm-5.3":             true,
	"kimi-k3":             true,
	"minimax-m3":          true,
}

func DefaultSettings(userID uuid.UUID) UserSettings {
	return UserSettings{
		UserID:           userID,
		DefaultRegion:    RegionSGP,
		DefaultModel:     "mimo-v2.6-flash",
		NotifyBalance:    true,
		NotifyUsage:      true,
		NotifySecurity:   true,
		NotifyProduct:    false,
		StreamDefault:    true,
		StoreRequestLogs: true,
		LoginAlerts:      true,
	}
}

type MagicPurpose string

const (
	PurposeLogin    MagicPurpose = "login"
	PurposeRegister MagicPurpose = "register"
)

type LoginAudit struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}
