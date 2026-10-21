package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/usecase"
)

const (
	ContextUserID    = "userID"
	ContextSessionID = "sessionID"
	ContextUser      = "user"
)

type SessionAuth struct {
	auth       *usecase.AuthUsecase
	cookieName string
	sessionTTL int // days
}

func NewSessionAuth(auth *usecase.AuthUsecase, cookieName string, sessionTTLDays int) *SessionAuth {
	return &SessionAuth{auth: auth, cookieName: cookieName, sessionTTL: sessionTTLDays}
}

func (m *SessionAuth) CookieName() string { return m.cookieName }

func (m *SessionAuth) CookieMaxAge() int { return m.sessionTTL * 24 * 3600 }

func (m *SessionAuth) RequireSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := c.Cookie(m.cookieName)
		if err != nil || raw == "" {
			// No cookie at all is normal for a signed-out visitor.
			LoggerFrom(c).Debug("session rejected: no cookie")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthenticated", "message": "not signed in"}})
			return
		}
		user, sess, err := m.auth.SessionUser(c.Request.Context(), raw)
		if err != nil {
			// A cookie was presented but did not resolve: either it expired or
			// it is bogus. Logged because a spike means something is wrong.
			LoggerFrom(c).Info("session rejected: lookup failed", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthenticated", "message": "not signed in"}})
			return
		}
		c.Set(ContextUser, user)
		c.Set(ContextUserID, user.ID)
		c.Set(ContextSessionID, sess.ID)
		c.Next()
	}
}

// RequireCSRF enforces a simple header check for cookie-authenticated mutations.
func RequireCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		if c.GetHeader("X-Requested-With") == "" && c.GetHeader("X-CSRF-Protection") == "" {
			LoggerFrom(c).Info("request rejected: missing CSRF header")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "csrf", "message": "missing CSRF header"}})
			return
		}
		c.Next()
	}
}

func ClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	return strings.TrimSpace(ip)
}

func UserAgent(c *gin.Context) string {
	return c.Request.UserAgent()
}

func UserFrom(c *gin.Context) (*entity.User, bool) {
	v, ok := c.Get(ContextUser)
	if !ok {
		return nil, false
	}
	u, ok := v.(*entity.User)
	return u, ok
}
