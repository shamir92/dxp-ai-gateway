package newapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	createUserPath = "/api/user/"
	maxDrainBytes  = 64 * 1024
)

type Client struct {
	baseURL    string
	adminToken string
	http       *http.Client
}

type APIError struct {
	StatusCode int
	RetryAfter time.Duration
	Retryable  bool
}

func (e *APIError) Error() string {
	return fmt.Sprintf("external API returned HTTP %d", e.StatusCode)
}

func NewClient(baseURL, adminToken string, timeout time.Duration) (*Client, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed == nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, fmt.Errorf("external API base URL must be an HTTPS origin")
	}
	if parsed.Path != "" && parsed.Path != "/" {
		return nil, fmt.Errorf("external API base URL must not include a path")
	}
	if strings.TrimSpace(adminToken) == "" {
		return nil, fmt.Errorf("external API admin token is required")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		baseURL:    baseURL,
		adminToken: adminToken,
		http: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

func (c *Client) CreateUser(ctx context.Context, username, displayName string, password []byte) error {
	payload := struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Password    string `json:"password"`
		Role        int    `json:"role"`
	}{Username: username, DisplayName: displayName, Password: string(password), Role: 1}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode external account request")
	}
	defer clear(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+createUserPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create external account request")
	}
	req.Header.Set("Authorization", "Bearer "+c.adminToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Referer", c.baseURL+"/users")
	req.Header.Set("Cookie", "new_api_has_session=1")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("external account request failed")
	}
	defer resp.Body.Close()
	_, _ = io.CopyN(io.Discard, resp.Body, maxDrainBytes)
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return nil
	}
	apiErr := &APIError{StatusCode: resp.StatusCode}
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= http.StatusInternalServerError {
		apiErr.Retryable = true
		apiErr.RetryAfter = parseRetryAfter(resp.Header.Get("Retry-After"), time.Now())
	}
	return apiErr
}

func parseRetryAfter(value string, now time.Time) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(value); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	if at, err := http.ParseTime(value); err == nil && at.After(now) {
		return at.Sub(now)
	}
	return 0
}
