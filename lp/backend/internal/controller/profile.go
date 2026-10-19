package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/middleware"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/oauth"
	"github.com/dxpglobal/ai-gateway/backend/internal/usecase"
)

type ProfileController struct {
	profile *usecase.ProfileUsecase
	auth    *usecase.AuthUsecase
	google  *oauth.GoogleProvider
	github  *oauth.GithubProvider
	apiKeys *usecase.APIKeyUsecase
}

func NewProfileController(profile *usecase.ProfileUsecase, auth *usecase.AuthUsecase, google *oauth.GoogleProvider, github *oauth.GithubProvider, apiKeys *usecase.APIKeyUsecase) *ProfileController {
	return &ProfileController{profile: profile, auth: auth, google: google, github: github, apiKeys: apiKeys}
}

func (h *ProfileController) Get(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	view, err := h.profile.GetProfile(c.Request.Context(), userID)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, view)
}

type patchProfile struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func (h *ProfileController) Patch(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	var req patchProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, usecase.ErrValidation)
		return
	}
	user, err := h.profile.UpdateProfile(c.Request.Context(), userID, req.Name, req.Email)
	if err != nil {
		writeErr(c, err)
		return
	}
	// Saving profile changes is also a chance to backfill a user whose external
	// accounts were never created. Best effort: the update already succeeded.
	if h.apiKeys != nil {
		if err := h.apiKeys.EnsureProvisioned(c.Request.Context(), userID); err != nil {
			log.Printf("profile: external account backfill failed for user %s", userID)
		}
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *ProfileController) ListSessions(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	sessionID := c.Value(middleware.ContextSessionID).(uuid.UUID)
	list, err := h.profile.ListSessions(c.Request.Context(), userID, sessionID)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"sessions": list})
}

func (h *ProfileController) RevokeSession(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeErr(c, usecase.ErrValidation)
		return
	}
	if err := h.profile.RevokeSession(c.Request.Context(), userID, id); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ProfileController) RevokeAllSessions(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	if err := h.profile.RevokeAllSessions(c.Request.Context(), userID); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ProfileController) ListConnections(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	list, err := h.profile.ListConnections(c.Request.Context(), userID)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"connections": list})
}

func (h *ProfileController) LinkStart(c *gin.Context) {
	provider := c.Param("provider")
	// Reuse OAuth start; callback links when already signed in.
	switch provider {
	case "google", "github":
		c.SetCookie("oauth_state", "st-link-"+provider, 600, "/", "", true, true)
		if provider == "google" {
			c.Redirect(http.StatusFound, h.google.AuthCodeURL("st-link-"+provider))
			return
		}
		c.Redirect(http.StatusFound, h.github.AuthCodeURL("st-link-"+provider))
	default:
		writeErr(c, usecase.ErrNotFound)
	}
}

func (h *ProfileController) Unlink(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	provider := entity.Provider(c.Param("provider"))
	if err := h.auth.UnlinkProvider(c.Request.Context(), userID, provider); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ProfileController) GetSecurity(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	view, err := h.profile.GetProfile(c.Request.Context(), userID)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"login_alerts": view.Settings.LoginAlerts})
}

type patchSecurity struct {
	LoginAlerts *bool `json:"login_alerts"`
}

func (h *ProfileController) PatchSecurity(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	var req patchSecurity
	if err := c.ShouldBindJSON(&req); err != nil || req.LoginAlerts == nil {
		writeErr(c, usecase.ErrValidation)
		return
	}
	s, err := h.profile.UpdateSecurity(c.Request.Context(), userID, *req.LoginAlerts)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"login_alerts": s.LoginAlerts})
}

type SettingsController struct {
	settings *usecase.SettingsUsecase
}

func NewSettingsController(settings *usecase.SettingsUsecase) *SettingsController {
	return &SettingsController{settings: settings}
}

func (h *SettingsController) Get(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	s, err := h.settings.Get(c.Request.Context(), userID)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *SettingsController) Patch(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	var req usecase.SettingsPatch
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, usecase.ErrValidation)
		return
	}
	s, err := h.settings.Patch(c.Request.Context(), userID, req)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, s)
}
