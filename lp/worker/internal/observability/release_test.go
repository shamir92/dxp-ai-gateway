package observability

import "testing"

// resolvedRelease caches, so these cases must call release() directly.
func TestReleasePriority(t *testing.T) {
	keys := []string{"SENTRY_RELEASE", "VERCEL_GIT_COMMIT_SHA", "GIT_COMMIT", "SOURCE_VERSION", "GITHUB_SHA"}
	clear := func() {
		for _, k := range keys {
			t.Setenv(k, "")
		}
	}

	t.Run("explicit override wins over everything", func(t *testing.T) {
		clear()
		t.Setenv("SENTRY_RELEASE", "v1.2.0")
		t.Setenv("VERCEL_GIT_COMMIT_SHA", "abc123")
		if got := release(); got != "v1.2.0" {
			t.Errorf("release() = %q, want v1.2.0", got)
		}
	})

	t.Run("falls back to the vercel commit sha", func(t *testing.T) {
		clear()
		t.Setenv("VERCEL_GIT_COMMIT_SHA", "abc123")
		if got := release(); got != "abc123" {
			t.Errorf("release() = %q, want abc123", got)
		}
	})

	t.Run("vercel sha wins over generic ci vars", func(t *testing.T) {
		clear()
		t.Setenv("VERCEL_GIT_COMMIT_SHA", "vercel-sha")
		t.Setenv("GITHUB_SHA", "github-sha")
		if got := release(); got != "vercel-sha" {
			t.Errorf("release() = %q, want vercel-sha", got)
		}
	})

	t.Run("generic git vars work outside vercel", func(t *testing.T) {
		clear()
		t.Setenv("SOURCE_VERSION", "src-1")
		if got := release(); got != "src-1" {
			t.Errorf("release() = %q, want src-1", got)
		}
	})

	t.Run("empty when nothing is set", func(t *testing.T) {
		clear()
		if got := release(); got != "" {
			t.Errorf("release() = %q, want empty", got)
		}
	})

	t.Run("whitespace is trimmed", func(t *testing.T) {
		clear()
		t.Setenv("VERCEL_GIT_COMMIT_SHA", "  abc  ")
		if got := release(); got != "abc" {
			t.Errorf("release() = %q, want abc", got)
		}
	})
}
