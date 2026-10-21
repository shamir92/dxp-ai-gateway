package observability

import (
	"testing"

	"github.com/getsentry/sentry-go"
)

func TestInitIsDisabledWithoutDSN(t *testing.T) {
	t.Setenv("SENTRY_DSN", "")
	if Init() {
		t.Fatal("Init reported enabled without a DSN")
	}
}

func TestScrubEventRemovesCredentialHeaders(t *testing.T) {
	event := &sentry.Event{
		Request: &sentry.Request{
			Headers: map[string]string{
				"Cookie":           "dxp_session=secret",
				"Authorization":    "Bearer nap_secret",
				"X-Security-Proof": "proof-secret",
				"Accept":           "application/json",
			},
			Cookies:     "dxp_session=secret",
			Data:        `{"password":"hunter2"}`,
			QueryString: "code=abc&access_token=nap_secret&page=2",
		},
		User: sentry.User{ID: "u1", Email: "someone@example.com"},
	}

	got := scrubEvent(event, nil)

	for _, header := range []string{"Cookie", "Authorization", "X-Security-Proof"} {
		if _, ok := got.Request.Headers[header]; ok {
			t.Errorf("%s header was not stripped", header)
		}
	}
	// Non-sensitive headers must survive so the event stays diagnosable.
	if got.Request.Headers["Accept"] != "application/json" {
		t.Errorf("Accept header was stripped")
	}
	if got.Request.Cookies != "" {
		t.Error("cookies were not cleared")
	}
	if got.Request.Data != "" {
		t.Error("request body was not cleared")
	}
	if got.User.Email != "" {
		t.Error("user email was not cleared")
	}
	if got.User.ID != "u1" {
		t.Error("user ID should be kept for grouping")
	}
	if want := "code=abc&access_token=[redacted]&page=2"; got.Request.QueryString != want {
		t.Errorf("QueryString = %q, want %q", got.Request.QueryString, want)
	}
}

func TestScrubBreadcrumbDropsSecretBodies(t *testing.T) {
	crumb := &sentry.Breadcrumb{
		Message: "POST /api/user/",
		Data: map[string]any{
			"body":          `{"password":"hunter2","token":"nap_secret"}`,
			"Authorization": "Bearer nap_secret",
			"status_code":   200,
		},
	}

	got := scrubBreadcrumb(crumb, nil)

	if _, ok := got.Data["body"]; ok {
		t.Error("body containing a secret was kept")
	}
	if _, ok := got.Data["Authorization"]; ok {
		t.Error("Authorization was kept")
	}
	if got.Data["status_code"] != 200 {
		t.Error("non-sensitive data was dropped")
	}
}

func TestRedactQuery(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", ""},
		{"page=2", "page=2"},
		{"token=abc", "token=[redacted]"},
		{"api_secret=abc", "api_secret=[redacted]"},
		{"code=nap_abc", "code=[redacted]"},
	} {
		if got := redactQuery(tc.in); got != tc.want {
			t.Errorf("redactQuery(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestTracesSampleRate(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want float64
	}{
		{"", 0.1},
		{"0.5", 0.5},
		{"0", 0},
		{"1", 1},
		{"2", 0.1},    // out of range falls back
		{"junk", 0.1}, // unparseable falls back
	} {
		t.Setenv("SENTRY_TRACES_SAMPLE_RATE", tc.raw)
		if got := tracesSampleRate(); got != tc.want {
			t.Errorf("tracesSampleRate(%q) = %v, want %v", tc.raw, got, tc.want)
		}
	}
}
