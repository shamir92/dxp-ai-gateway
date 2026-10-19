package newapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// TokenStatus values reported by the upstream.
const (
	TokenStatusEnabled   = 1
	TokenStatusDisabled  = 2
	TokenStatusExpired   = 3
	TokenStatusExhausted = 4
)

// Token is a new-api API key. The Key field is always masked in list responses.
type Token struct {
	ID                 int      `json:"id"`
	UserID             int      `json:"user_id"`
	Key                string   `json:"key"`
	Status             int      `json:"status"`
	Name               string   `json:"name"`
	CreatedTime        int64    `json:"created_time"`
	AccessedTime       int64    `json:"accessed_time"`
	ExpiredTime        int64    `json:"expired_time"`
	RemainQuota        int      `json:"remain_quota"`
	UnlimitedQuota     bool     `json:"unlimited_quota"`
	ModelLimitsEnabled bool     `json:"model_limits_enabled"`
	ModelLimits        string   `json:"model_limits"`
	AllowIPs           *string  `json:"allow_ips"`
	UsedQuota          int      `json:"used_quota"`
	Group              string   `json:"group"`
	CrossGroupRetry    bool     `json:"cross_group_retry"`
	AutoGroups         []string `json:"auto_groups"`
}

// TokenInput is the create/update payload for /api/token/. ID is omitted on
// create and required on update.
type TokenInput struct {
	ID                 int      `json:"id,omitempty"`
	Name               string   `json:"name"`
	RemainQuota        int      `json:"remain_quota"`
	ExpiredTime        int64    `json:"expired_time"`
	UnlimitedQuota     bool     `json:"unlimited_quota"`
	ModelLimitsEnabled bool     `json:"model_limits_enabled"`
	ModelLimits        string   `json:"model_limits"`
	AllowIPs           string   `json:"allow_ips"`
	Group              string   `json:"group"`
	AutoGroups         []string `json:"auto_groups"`
	CrossGroupRetry    bool     `json:"cross_group_retry"`
}

// TokenPage is the paged list payload. Items may be null upstream, so callers
// should treat a nil slice as empty.
type TokenPage struct {
	Items    []Token `json:"items"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Total    int     `json:"total"`
}

// ListTokens returns one page of the user's API keys. The upstream masks the key
// value, so no plaintext key is available here.
func (c *Client) ListTokens(ctx context.Context, bearer string, page, pageSize int) (TokenPage, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	query := url.Values{}
	query.Set("p", strconv.Itoa(page))
	query.Set("page_size", strconv.Itoa(pageSize))

	data, err := c.do(ctx, requestOptions{
		method:  http.MethodGet,
		path:    tokenPath + "?" + query.Encode(),
		bearer:  bearer,
		referer: "/keys",
	})
	if err != nil {
		return TokenPage{}, err
	}

	var result TokenPage
	if err := json.Unmarshal(data, &result); err != nil {
		return TokenPage{}, fmt.Errorf("decode token list response")
	}
	for i := range result.Items {
		if result.Items[i].AutoGroups == nil {
			result.Items[i].AutoGroups = []string{}
		}
	}
	return result, nil
}

// GetToken returns a single API key by id.
func (c *Client) GetToken(ctx context.Context, bearer string, id int) (Token, error) {
	data, err := c.do(ctx, requestOptions{
		method:  http.MethodGet,
		path:    tokenPath + strconv.Itoa(id),
		bearer:  bearer,
		referer: "/keys",
	})
	if err != nil {
		return Token{}, err
	}

	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return Token{}, fmt.Errorf("decode token response")
	}
	return token, nil
}

// CreateToken creates an API key. The plaintext value is only obtainable through
// RevealTokenKey.
func (c *Client) CreateToken(ctx context.Context, bearer string, input TokenInput) error {
	if input.AutoGroups == nil {
		input.AutoGroups = []string{}
	}
	input.ID = 0
	_, err := c.do(ctx, requestOptions{
		method:  http.MethodPost,
		path:    tokenPath,
		bearer:  bearer,
		referer: "/keys",
		body:    input,
	})
	return err
}

// UpdateToken updates an existing API key. The id travels in the body.
func (c *Client) UpdateToken(ctx context.Context, bearer string, input TokenInput) (Token, error) {
	if input.AutoGroups == nil {
		input.AutoGroups = []string{}
	}
	data, err := c.do(ctx, requestOptions{
		method:  http.MethodPut,
		path:    tokenPath,
		bearer:  bearer,
		referer: "/keys",
		body:    input,
	})
	if err != nil {
		return Token{}, err
	}

	var token Token
	if len(data) > 0 {
		if err := json.Unmarshal(data, &token); err != nil {
			return Token{}, fmt.Errorf("decode token response")
		}
	}
	return token, nil
}

// DeleteToken removes an API key by id.
func (c *Client) DeleteToken(ctx context.Context, bearer string, id int) error {
	_, err := c.do(ctx, requestOptions{
		method:  http.MethodDelete,
		path:    tokenPath + strconv.Itoa(id),
		bearer:  bearer,
		referer: "/keys",
	})
	return err
}

// RevealTokenKey returns the plaintext key for one API key. It is the only read
// path that exposes a full key, and callers must not persist the result.
func (c *Client) RevealTokenKey(ctx context.Context, bearer string, id int) (string, error) {
	data, err := c.do(ctx, requestOptions{
		method:  http.MethodPost,
		path:    tokenPath + strconv.Itoa(id) + "/key",
		bearer:  bearer,
		referer: "/keys",
	})
	if err != nil {
		return "", err
	}

	var result struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("decode token key response")
	}
	if strings.TrimSpace(result.Key) == "" {
		return "", fmt.Errorf("token key response did not include a key")
	}
	return result.Key, nil
}
