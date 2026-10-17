package entity

import (
	"time"

	"github.com/google/uuid"
)

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
