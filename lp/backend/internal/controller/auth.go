package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/middleware"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/oauth"
	"github.com/dxpglobal/ai-gateway/backend/internal/usecase"
)

// writeErr maps a domain error to a client response.
//
// The sentinel 4xx cases are expected traffic and stay quiet: a wrong code or
// an expired session is not an incident. Unexpected failures (and upstream
// faults, which mean the external platform is broken rather than the caller)
// are reported, because previously the 500 path answered "internal error"
// without recording the cause anywhere.
func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrValidation):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "validation", "message": err.Error()}})
	case errors.Is(err, usecase.ErrUnauthenticated):
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthenticated", "message": "not signed in"}})
	case errors.Is(err, usecase.ErrLastAuthMethod):
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "last_auth_method", "message": "cannot remove your only login method"}})
	case errors.Is(err, usecase.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "not found"}})
	case errors.Is(err, usecase.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "conflict", "message": "already linked to another account"}})
	case errors.Is(err, usecase.ErrRateLimited), errors.Is(err, usecase.ErrTooManyAttempts):
		c.JSON(http.StatusTooManyRequests, gin.H{"error": gin.H{"code": "rate_limited", "message": "too many requests"}})
	case errors.Is(err, usecase.ErrCodeExpired), errors.Is(err, usecase.ErrInvalidCredentials):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "validation", "message": "invalid or expired code"}})
	case errors.Is(err, usecase.ErrNoAccount):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "no_account", "message": "No account for this email. Create one first."}})
	case errors.Is(err, usecase.ErrAccountExists):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "account_exists", "message": "Account already exists. Log in instead."}})
	case errors.Is(err, usecase.ErrUnsupportedSite):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "site_unsupported", "message": "unsupported external site"}})
	case errors.Is(err, usecase.ErrExternalAccountMissing):
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "external_account_missing", "message": "your account is still being provisioned on this platform"}})
	case errors.Is(err, usecase.ErrExternalCredentials):
		reportUpstream(c, err, "external platform rejected our credentials")
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": "external_credentials", "message": "the external platform rejected our credentials"}})
	case errors.Is(err, usecase.ErrExternalRejected):
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "external_rejected", "message": "the external platform rejected the request"}})
	case errors.Is(err, usecase.ErrExternalUnavailable):
		reportUpstream(c, err, "external platform unavailable")
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": "external_unavailable", "message": "the external platform is unavailable"}})
	default:
		reportInternal(c, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal", "message": "internal error"}})
	}
}

// reportInternal logs an unexpected fault and sends it to Sentry. This is the
// path that previously answered 500 with no trace of the underlying error.
func reportInternal(c *gin.Context, err error) {
	log := middleware.LoggerFrom(c).With(
		slog.String("error_code", "internal"),
		slog.String("error", err.Error()),
	)
	log.Error("unhandled request error")
	sentry.CaptureException(err)
}

// reportUpstream records a failed call to an external platform. These are 5xx
// because the fault is upstream, so they are worth surfacing.
func reportUpstream(c *gin.Context, err error, reason string) {
	log := middleware.LoggerFrom(c).With(
		slog.String("error_code", "upstream"),
		slog.String("reason", reason),
		slog.String("error", err.Error()),
	)
	log.Error("upstream request failed")
	sentry.CaptureException(err)
}

type AuthController struct {
	auth         *usecase.AuthUsecase
	sessionsMW   *middleware.SessionAuth
	google       *oauth.GoogleProvider
	github       *oauth.GithubProvider
	publicWebURL string
	stateStore   map[string]string
}

func NewAuthController(
	auth *usecase.AuthUsecase,
	sessionsMW *middleware.SessionAuth,
	google *oauth.GoogleProvider,
	github *oauth.GithubProvider,
	publicWebURL string,
) *AuthController {
	return &AuthController{
		auth:         auth,
		sessionsMW:   sessionsMW,
		google:       google,
		github:       github,
		publicWebURL: publicWebURL,
		stateStore:   map[string]string{},
	}
}

func (h *AuthController) setSessionCookie(c *gin.Context, token string) {
	c.SetCookie(
		h.sessionsMW.CookieName(),
		token,
		h.sessionsMW.CookieMaxAge(),
		"/",
		"",
		true,
		true,
	)
}

func (h *AuthController) clearSessionCookie(c *gin.Context) {
	c.SetCookie(h.sessionsMW.CookieName(), "", -1, "/", "", true, true)
}

type magicRequest struct {
	Email   string `json:"email" binding:"required"`
	Purpose string `json:"purpose" binding:"required"`
	Name    string `json:"name"`
}

func (h *AuthController) MagicRequest(c *gin.Context) {
	var req magicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, usecase.ErrValidation)
		return
	}
	// 400 when the email does not match the flow (login vs register).
	err := h.auth.RequestMagicCode(c.Request.Context(), req.Email, req.Purpose, req.Name)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *AuthController) MagicResend(c *gin.Context) {
	var req magicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, usecase.ErrValidation)
		return
	}
	err := h.auth.ResendMagicCode(c.Request.Context(), req.Email, req.Purpose)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type magicVerify struct {
	Email   string `json:"email" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Purpose string `json:"purpose" binding:"required"`
}

func (h *AuthController) MagicVerify(c *gin.Context) {
	var req magicVerify
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, usecase.ErrValidation)
		return
	}
	meta := usecase.SessionMeta{UserAgent: middleware.UserAgent(c), IP: middleware.ClientIP(c)}
	res, err := h.auth.VerifyMagicCode(c.Request.Context(), req.Email, req.Code, entity.MagicPurpose(req.Purpose), meta)
	if err != nil {
		writeErr(c, err)
		return
	}
	h.setSessionCookie(c, res.Token)
	c.JSON(http.StatusOK, gin.H{"user": res.User})
}

func (h *AuthController) OAuthStart(c *gin.Context) {
	provider := c.Param("provider")
	state := "st-" + provider // real state should be random+cookie; keep simple + cookie below
	c.SetCookie("oauth_state", state, 600, "/", "", true, true)
	var url string
	switch provider {
	case "google":
		url = h.google.AuthCodeURL(state)
	case "github":
		url = h.github.AuthCodeURL(state)
	default:
		writeErr(c, usecase.ErrNotFound)
		return
	}
	c.Redirect(http.StatusFound, url)
}

func (h *AuthController) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	want, _ := c.Cookie("oauth_state")
	if state == "" || want == "" || state != want {
		c.Redirect(http.StatusFound, h.publicWebURL+"/login?error=oauth")
		return
	}
	var prof *oauth.Profile
	var err error
	switch provider {
	case "google":
		prof, err = h.google.Exchange(c.Request.Context(), code)
	case "github":
		prof, err = h.github.Exchange(c.Request.Context(), code)
	default:
		writeErr(c, usecase.ErrNotFound)
		return
	}
	if err != nil {
		c.Redirect(http.StatusFound, h.publicWebURL+"/login?error=oauth")
		return
	}
	meta := usecase.SessionMeta{UserAgent: middleware.UserAgent(c), IP: middleware.ClientIP(c)}
	res, err := h.auth.CompleteOAuthLogin(
		c.Request.Context(),
		entity.Provider(provider),
		prof.ProviderUserID,
		prof.Email,
		prof.Name,
		prof.AvatarURL,
		meta,
	)
	if err != nil {
		c.Redirect(http.StatusFound, h.publicWebURL+"/login?error=oauth")
		return
	}
	h.setSessionCookie(c, res.Token)
	c.Redirect(http.StatusFound, h.publicWebURL+"/token-plan")
}

func (h *AuthController) Logout(c *gin.Context) {
	if raw, ok := c.Get(middleware.ContextSessionID); ok {
		if id, ok := raw.(uuid.UUID); ok {
			_ = h.auth.Logout(c.Request.Context(), id)
		}
	}
	h.clearSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *AuthController) Me(c *gin.Context) {
	user, ok := middleware.UserFrom(c)
	if !ok {
		writeErr(c, usecase.ErrUnauthenticated)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
