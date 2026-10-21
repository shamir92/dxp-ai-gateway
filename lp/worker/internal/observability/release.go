package observability

import (
	"os"
	"strings"
	"sync"
)

// release resolves the value reported to Sentry as the release.
//
// Priority:
//  1. SENTRY_RELEASE - an explicit override always wins.
//  2. VERCEL_GIT_COMMIT_SHA - Vercel injects the deployed commit, available at
//     both build and runtime (requires "Enable access to System Environment
//     Variables" in project settings).
//  3. GIT_COMMIT / SOURCE_VERSION / GITHUB_SHA - the same idea elsewhere.
//
// Vercel exposes no tag variable, so a tag like "v1.2.0" cannot be read
// directly. The deploy workflow passes it through SENTRY_RELEASE at build time,
// which is why that takes priority.
func release() string {
	for _, key := range []string{
		"SENTRY_RELEASE",
		"VERCEL_GIT_COMMIT_SHA",
		"GIT_COMMIT",
		"SOURCE_VERSION",
		"GITHUB_SHA",
	} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	return ""
}

// resolvedRelease caches the answer: the environment does not change while the
// process runs, and release() is called on every Init.
var (
	releaseOnce  sync.Once
	releaseCache string
)

func resolvedRelease() string {
	releaseOnce.Do(func() { releaseCache = release() })
	return releaseCache
}
