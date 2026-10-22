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
	createUserPath      = "/api/user/"
	maxDrainBytes       = 64 * 1024
	sessionMarkerCookie = "new_api_has_session=1"
	// Cloudflare Access service-token headers. Without these the Zero Trust
	// policy rejects the request before it ever reaches new-api.
	cfAccessClientIDHeader     = "CF-Access-Client-Id"
	cfAccessClientSecretHeader = "CF-Access-Client-Secret"
)

// Client talks to a new-api instance with the site's admin bearer. The worker
// only ever creates users; token, key and subscription calls stay in the backend.
type Client struct {
	baseURL        string
	adminToken     string
	cfAccessID     string
	cfAccessSecret string
	http           *http.Client
}

type APIError struct {
	StatusCode int
	Code       string
	Message    string
	RetryAfter time.Duration
	Retryable  bool
}

func (e *APIError) Error() string {
	// Include the upstream code/message: without it a rejected call is
	// indistinguishable from any other 4xx, which makes diagnosis impossible.
	detail := strings.TrimSpace(e.Message)
	if detail == "" {
		detail = strings.TrimSpace(e.Code)
	}
	if detail == "" {
		return fmt.Sprintf("external API returned HTTP %d", e.StatusCode)
	}
	if e.Code != "" && e.Message != "" {
		return fmt.Sprintf("external API returned HTTP %d (%s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("external API returned HTTP %d: %s", e.StatusCode, detail)
}

// NewClient builds a client for one new-api site. cfAccessID/cfAccessSecret are
// the Cloudflare Access service token; both or neither must be set.
func NewClient(baseURL, adminToken, cfAccessID, cfAccessSecret string, timeout time.Duration) (*Client, error) {
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
	cfAccessID = strings.TrimSpace(cfAccessID)
	cfAccessSecret = strings.TrimSpace(cfAccessSecret)
	if (cfAccessID == "") != (cfAccessSecret == "") {
		return nil, fmt.Errorf("Cloudflare Access client ID and secret must both be set")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		baseURL:        baseURL,
		adminToken:     adminToken,
		cfAccessID:     cfAccessID,
		cfAccessSecret: cfAccessSecret,
		http: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

// envelope is the new-api response wrapper. Branch on Success, never on the HTTP
// status: business errors are returned with HTTP 200 and success=false.
type envelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Code    string          `json:"code"`
	Data    json.RawMessage `json:"data"`
}

type requestOptions struct {
	method  string
	path    string
	bearer  string
	referer string
	body    any
	headers map[string]string
}

// do performs one JSON request and returns the raw `data` field. It never logs or
// returns the request body, so credentials/keys cannot leak through it.
func (c *Client) do(ctx context.Context, opts requestOptions) (json.RawMessage, error) {
	var reader io.Reader
	if opts.body != nil {
		payload, err := json.Marshal(opts.body)
		if err != nil {
			return nil, fmt.Errorf("encode external request")
		}
		defer clear(payload)
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, opts.method, c.baseURL+opts.path, reader)
	if err != nil {
		return nil, fmt.Errorf("create external request")
	}
	if opts.bearer != "" {
		req.Header.Set("Authorization", "Bearer "+opts.bearer)
	}
	req.Header.Set("Accept", "application/json")
	if opts.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Origin", c.baseURL)
	referer := opts.referer
	if referer == "" {
		referer = "/"
	}
	req.Header.Set("Referer", c.baseURL+referer)
	req.Header.Set("Cookie", sessionMarkerCookie)
	// Cloudflare Access sits in front of the site, so the service token is
	// required on every call.
	if c.cfAccessID != "" {
		req.Header.Set(cfAccessClientIDHeader, c.cfAccessID)
		req.Header.Set(cfAccessClientSecretHeader, c.cfAccessSecret)
	}
	for name, value := range opts.headers {
		req.Header.Set(name, value)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("external request failed")
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxDrainBytes))
	if err != nil {
		return nil, fmt.Errorf("read external response")
	}

	apiError := func(code, message string) *APIError {
		e := &APIError{StatusCode: resp.StatusCode, Code: code, Message: message}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= http.StatusInternalServerError {
			e.Retryable = true
			e.RetryAfter = parseRetryAfter(resp.Header.Get("Retry-After"), time.Now())
		}
		return e
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, apiError("", "")
	}

	var env envelope
	if len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, &env); err != nil {
			return nil, fmt.Errorf("decode external response")
		}
	}
	if !env.Success {
		return nil, apiError(env.Code, env.Message)
	}
	return env.Data, nil
}

// CreateUser provisions a remote user with the site admin bearer.
//
// new-api reports business failures as HTTP 200 with success=false, so the
// response envelope MUST be checked. Treating any 2xx as success marked accounts
// provisioned that were never created — which later surfaced as a login failure.
func (c *Client) CreateUser(ctx context.Context, username, displayName string, password []byte) error {
	payload := struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Password    string `json:"password"`
		Role        int    `json:"role"`
	}{Username: username, DisplayName: displayName, Password: string(password), Role: 1}

	_, err := c.do(ctx, requestOptions{
		method:  http.MethodPost,
		path:    createUserPath,
		bearer:  c.adminToken,
		referer: "/users",
		body:    payload,
	})
	return err
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
