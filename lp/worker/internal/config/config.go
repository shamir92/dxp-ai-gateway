package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// Config is the subset of settings the provisioning worker needs. It shares the
// same environment variable names as the backend so both can read one .env.
type Config struct {
	DatabaseURL                   string
	TPAI                          ExternalSite
	APAI                          ExternalSite
	ExternalPasswordEncryptionKey []byte
}

type ExternalSite struct {
	BaseURL    string
	AdminToken string
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
