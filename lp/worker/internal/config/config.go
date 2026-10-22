package config

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Config is the subset of settings the provisioning worker needs. The names
// match the backend's so the same values can be copied between the two files,
// but the worker reads its own lp/worker/.env and shares no file with the
// backend.
type Config struct {
	DatabaseURL                   string
	TPAI                          ExternalSite
	APAI                          ExternalSite
	ExternalPasswordEncryptionKey []byte
	// CloudflareAccess guards the external sites behind Zero Trust, so every
	// request to TPAI/APAI must carry a service token.
	CloudflareAccess CloudflareAccess
}

type ExternalSite struct {
	BaseURL    string
	AdminToken string
}

type CloudflareAccess struct {
	ClientID     string
	ClientSecret string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		TPAI: ExternalSite{
			BaseURL:    strings.TrimRight(strings.TrimSpace(os.Getenv("TPAI_BASE_URL")), "/"),
			AdminToken: strings.TrimSpace(os.Getenv("TPAI_ADMIN_TOKEN")),
		},
		APAI: ExternalSite{
			BaseURL:    strings.TrimRight(strings.TrimSpace(os.Getenv("APAI_BASE_URL")), "/"),
			AdminToken: strings.TrimSpace(os.Getenv("APAI_ADMIN_TOKEN")),
		},
		CloudflareAccess: CloudflareAccess{
			ClientID:     strings.TrimSpace(os.Getenv("CF_CLIENT_ID")),
			ClientSecret: strings.TrimSpace(os.Getenv("CF_CLIENT_SECRET")),
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
	if (cfg.CloudflareAccess.ClientID == "") != (cfg.CloudflareAccess.ClientSecret == "") {
		return nil, fmt.Errorf("CF_CLIENT_ID and CF_CLIENT_SECRET must both be set")
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

// DatabaseTarget renders a DSN as "host:port/database" for logging. Credentials
// are dropped, so the result is safe to emit: the worker's logs must never carry
// the database password.
func DatabaseTarget(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil || u.Host == "" {
		return "unknown"
	}
	return u.Host + u.Path
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
