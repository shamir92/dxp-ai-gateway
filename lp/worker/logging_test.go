package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

// The worker's output must be machine-readable: it runs outside the platform's
// request pipeline, so there is no request log to fall back on.
func TestWorkerLogsAreJSON(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("claim failed", slog.Any("error", errBoom{}))

	os.Stdout = orig
	_ = w.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	line := strings.TrimSpace(string(data))
	if !strings.HasPrefix(line, "{") {
		t.Fatalf("not JSON: %q", line)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("invalid JSON %q: %v", line, err)
	}
	if got["msg"] != "claim failed" {
		t.Errorf("msg = %v", got["msg"])
	}
	// The whole point of the change: the error must actually be present.
	if _, ok := got["error"]; !ok {
		t.Error("error field is missing from the log line")
	}
}

type errBoom struct{}

func (errBoom) Error() string { return "connection refused" }
