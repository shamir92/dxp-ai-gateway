package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr                      string
	DatabaseURL                   string
	SessionCookieName             string
	SessionTTLDays                int
	GoogleClientID                string
	GoogleClientSecret            string
	GoogleRedirectURL             string
	GithubClientID                string
	GithubClientSecret            string
	GithubRedirectURL             string
	SMTPHost                      string
	SMTPPort                      int
	SMTPUsername                  string
	SMTPPassword                  string
	SMTPFrom                      string
	SMTPFromName                  string
	SMTPTLS                       string // starttls | tls | none
	MailBackend                   string // smtp | console
	PublicWebURL                  string
	TPAI                          ExternalSite
	APAI                          ExternalSite
	ExternalPasswordEncryptionKey []byte
}

type ExternalSite struct {
	BaseURL    string
	AdminToken string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getenv("SMTP_PORT", "587"))
	if err != nil {
		return nil, fmt.Errorf("SMTP_PORT: %w", err)
	}
	ttl, err := strconv.Atoi(getenv("SESSION_TTL_DAYS", "30"))
	if err != nil {
		return nil, fmt.Errorf("SESSION_TTL_DAYS: %w", err)
	}
	cfg := &Config{
		HTTPAddr:           httpAddr(),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		SessionCookieName:  getenv("SESSION_COOKIE_NAME", "dxp_session"),
		SessionTTLDays:     ttl,
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubRedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		SMTPHost:           os.Getenv("SMTP_HOST"),
		SMTPPort:           port,
		SMTPUsername:       os.Getenv("SMTP_USERNAME"),
		SMTPPassword:       os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:           getenv("SMTP_FROM", "hello@fromdxp.com"),
		SMTPFromName:       getenv("SMTP_FROM_NAME", "DXP AI Gateway"),
		SMTPTLS:            getenv("SMTP_TLS", "starttls"),
		MailBackend:        getenv("MAIL_BACKEND", "smtp"),
		PublicWebURL:       getenv("PUBLIC_WEB_URL", "http://localhost:3000"),
		TPAI: ExternalSite{
			BaseURL:    strings.TrimRight(strings.TrimSpace(os.Getenv("TPAI_BASE_URL")), "/"),
			AdminToken: strings.TrimSpace(os.Getenv("TPAI_ADMIN_TOKEN")),
		},
		APAI: ExternalSite{
			BaseURL:    strings.TrimRight(strings.TrimSpace(os.Getenv("APAI_BASE_URL")), "/"),
			AdminToken: strings.TrimSpace(os.Getenv("APAI_ADMIN_TOKEN")),
		},
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	for name, site := range map[string]ExternalSite{"TPAI": cfg.TPAI, "APAI": cfg.APAI} {
		if site.BaseURL == "" || site.AdminToken == "" {
			return nil, fmt.Errorf("%s_BASE_URL and %s_ADMIN_TOKEN are required", name, name)
		}
	}
	keyText := strings.TrimSpace(os.Getenv("EXTERNAL_PASSWORD_ENCRYPTION_KEY"))
	if keyText == "" {
		return nil, fmt.Errorf("EXTERNAL_PASSWORD_ENCRYPTION_KEY is required")
	}
	key, err := base64.StdEncoding.DecodeString(keyText)
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("EXTERNAL_PASSWORD_ENCRYPTION_KEY must be base64-encoded 32-byte key")
	}
	cfg.ExternalPasswordEncryptionKey = key
	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// httpAddr prefers the platform-assigned PORT (Vercel) and falls back to
// HTTP_ADDR, then the local default.
func httpAddr() string {
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		return ":" + p
	}
	return getenv("HTTP_ADDR", ":8080")
}

// LoadDotEnv reads KEY=VALUE lines into the process environment if the key
// is not already set. Missing file is fine. Comments (#) and blank lines ignored.
func LoadDotEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		v = strings.Trim(v, `"'`)
		if k == "" {
			continue
		}
		if _, exists := os.LookupEnv(k); !exists {
			_ = os.Setenv(k, v)
		}
	}
	return nil
}
