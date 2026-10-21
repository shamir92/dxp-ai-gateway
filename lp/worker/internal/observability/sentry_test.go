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

// The worker holds decrypted site passwords, so a body reaching an event is a
// credential leak. Header names arrive canonical-case from Sentry ("Cookie"),
// which a lowercase-only comparison would silently miss.
func TestScrubEventRemovesCredentials(t *testing.T) {
	event := &sentry.Event{
		Request: &sentry.Request{
			Headers: map[string]string{
				"Cookie":        "dxp_session=secret",
				"Authorization": "Bearer nap_secret",
				"Accept":        "application/json",
			},
			Cookies: "dxp_session=secret",
			Data:    `{"password":"hunter2"}`,
		},
		User: sentry.User{ID: "u1", Email: "someone@example.com"},
	}

	got := scrubEvent(event, nil)

	if _, ok := got.Request.Headers["Cookie"]; ok {
		t.Error("Cookie header was not stripped")
	}
	if _, ok := got.Request.Headers["Authorization"]; ok {
		t.Error("Authorization header was not stripped")
	}
	if got.Request.Headers["Accept"] != "application/json" {
		t.Error("non-sensitive header was stripped")
	}
	if got.Request.Data != "" {
		t.Error("request body was not cleared")
	}
	if got.Request.Cookies != "" {
		t.Error("cookies were not cleared")
	}
	if got.User.Email != "" {
		t.Error("user email was not cleared")
	}
}
