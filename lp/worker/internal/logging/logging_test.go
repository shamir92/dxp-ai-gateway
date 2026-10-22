package logging

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

// capture runs fn with os.Stdout redirected and returns what was written.
func capture(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return string(data)
}

// A bare slog call must produce the same JSON as the package logger. This is the
// regression the backend hit: the package logger was JSON while the stdlib
// default handler stayed text, so startup lines came out as `2006/01/02 INFO …`.
// The worker had exactly that gap before this package existed.
func TestSetupMakesDefaultSlogEmitJSON(t *testing.T) {
	out := capture(t, func() {
		Setup()
		slog.Info("startup line", slog.String("stage", "test"))
	})

	line := strings.TrimSpace(out)
	if !strings.HasPrefix(line, "{") {
		t.Fatalf("default slog is not JSON: %q", line)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("invalid JSON %q: %v", line, err)
	}
	if got["msg"] != "startup line" {
		t.Errorf("msg = %v, want %q", got["msg"], "startup line")
	}
	if got["stage"] != "test" {
		t.Errorf("stage = %v, want %q", got["stage"], "test")
	}
}

func TestLevelFromEnv(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"", slog.LevelInfo},
		{"nonsense", slog.LevelInfo},
	} {
		if got := parseLevel(tc.raw); got != tc.want {
			t.Errorf("parseLevel(%q) = %v, want %v", tc.raw, got, tc.want)
		}
	}
}
