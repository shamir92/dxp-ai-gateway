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
