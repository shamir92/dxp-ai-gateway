package entity

import (
	"time"

	"github.com/google/uuid"
)

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
