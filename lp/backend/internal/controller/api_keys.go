package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/newapi"
	"github.com/dxpglobal/ai-gateway/backend/internal/middleware"
	"github.com/dxpglobal/ai-gateway/backend/internal/usecase"
)

// maxKeyNameLength mirrors the upstream limit for a token name.
const maxKeyNameLength = 50

// contextSite carries the backing platform for a route whose public path uses a
// neutral product slug (tp_, ap_) so the upstream platform is never named in a
// URL or response.
const contextSite = "external_site"

// WithSite binds the backing platform to a route.
func WithSite(site entity.ExternalSite) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(contextSite, site)
		c.Next()
	}
}

func siteFrom(c *gin.Context) entity.ExternalSite {
	return c.MustGet(contextSite).(entity.ExternalSite)
}

type APIKeyController struct {
	keys *usecase.APIKeyUsecase
}

func NewAPIKeyController(keys *usecase.APIKeyUsecase) *APIKeyController {
	return &APIKeyController{keys: keys}
}

// keyInputBody is the create/update payload. Every field is optional except the
// name, so callers send only what they care about and upstream-compatible
// defaults fill the rest.
type keyInputBody struct {
	Name               *string  `json:"name"`
	UnlimitedQuota     *bool    `json:"unlimited_quota"`
	ExpiredTime        *int64   `json:"expired_time"`
	RemainQuota        *int     `json:"remain_quota"`
	ModelLimitsEnabled *bool    `json:"model_limits_enabled"`
	ModelLimits        *string  `json:"model_limits"`
	AllowIPs           *string  `json:"allow_ips"`
	Group              *string  `json:"group"`
	AutoGroups         []string `json:"auto_groups"`
	CrossGroupRetry    *bool    `json:"cross_group_retry"`
}

// toInput validates the name and applies defaults matching the upstream form.
func (b keyInputBody) toInput(id int) (newapi.TokenInput, error) {
	input := newapi.TokenInput{
		ID:          id,
		AutoGroups:  b.AutoGroups,
		ExpiredTime: -1,
	}
	if input.AutoGroups == nil {
		input.AutoGroups = []string{}
	}
	if b.Name == nil {
		return newapi.TokenInput{}, usecase.ErrValidation
	}
	name := strings.TrimSpace(*b.Name)
	if name == "" || len([]rune(name)) > maxKeyNameLength {
		return newapi.TokenInput{}, usecase.ErrValidation
	}
	input.Name = name

	if b.UnlimitedQuota != nil {
		input.UnlimitedQuota = *b.UnlimitedQuota
	}
	if b.ExpiredTime != nil {
		input.ExpiredTime = *b.ExpiredTime
	}
	if b.RemainQuota != nil {
		if *b.RemainQuota < 0 {
			return newapi.TokenInput{}, usecase.ErrValidation
		}
		input.RemainQuota = *b.RemainQuota
	}
	if b.ModelLimitsEnabled != nil {
		input.ModelLimitsEnabled = *b.ModelLimitsEnabled
	}
	if b.ModelLimits != nil {
		input.ModelLimits = *b.ModelLimits
	}
	if b.AllowIPs != nil {
		input.AllowIPs = *b.AllowIPs
	}
	if b.Group != nil {
		input.Group = *b.Group
	}
	if b.CrossGroupRetry != nil {
		input.CrossGroupRetry = *b.CrossGroupRetry
	}
	return input, nil
}

func (h *APIKeyController) List(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	// Opening the plan page checks the user's external accounts and enqueues any
	// missing ones, so a user who registered before provisioning existed is healed
	// here instead of getting a permanent "not provisioned" error.
	if err := h.keys.EnsureProvisioned(c.Request.Context(), userID); err != nil {
		log.Printf("api keys: external account backfill failed user=%s", userID)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.keys.ListKeys(c.Request.Context(), userID, siteFrom(c), page, pageSize)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":     result.Items,
		"page":      result.Page,
		"page_size": result.PageSize,
		"total":     result.Total,
	})
}

func (h *APIKeyController) Create(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	var body keyInputBody
	if err := c.ShouldBindJSON(&body); err != nil {
		writeErr(c, usecase.ErrValidation)
		return
	}
	input, err := body.toInput(0)
	if err != nil {
		writeErr(c, err)
		return
	}

	created, err := h.keys.CreateKey(c.Request.Context(), userID, siteFrom(c), input)
	if err != nil {
		writeErr(c, err)
		return
	}
	// Key is the plaintext value. The upstream exposes it once, here, and never
	// again through the list endpoint.
	c.JSON(http.StatusOK, gin.H{
		"token": created.Token,
		"key":   created.Key,
	})
}

func (h *APIKeyController) Update(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		writeErr(c, usecase.ErrValidation)
		return
	}
	var body keyInputBody
	if err := c.ShouldBindJSON(&body); err != nil {
		writeErr(c, usecase.ErrValidation)
		return
	}
	input, err := body.toInput(id)
	if err != nil {
		writeErr(c, err)
		return
	}

	updated, err := h.keys.UpdateKey(c.Request.Context(), userID, siteFrom(c), input)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": updated})
}

func (h *APIKeyController) Delete(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		writeErr(c, usecase.ErrValidation)
		return
	}

	if err := h.keys.DeleteKey(c.Request.Context(), userID, siteFrom(c), id); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *APIKeyController) Subscription(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	view, err := h.keys.Subscription(c.Request.Context(), userID, siteFrom(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, view)
}

// Endpoints returns the platform's advertised API base addresses. It is public
// upstream, so it needs no stored credential and cannot fail on a missing account.
func (h *APIKeyController) Endpoints(c *gin.Context) {
	userID := c.Value(middleware.ContextUserID).(uuid.UUID)
	list, err := h.keys.Endpoints(c.Request.Context(), userID, siteFrom(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}
