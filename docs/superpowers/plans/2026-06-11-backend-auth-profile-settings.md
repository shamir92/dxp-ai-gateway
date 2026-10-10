# Backend M1 — Auth, Profile, Settings Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Go/Gin control-plane API for magic-code + Google/GitHub auth, profile, settings, logout.

**Architecture:** Clean architecture `controller → usecase → repository → entity` in `lp/backend`. Opaque sessions in Postgres (hashed token in httpOnly cookie). SMTP magic-code mail with backend-owned Go templates.

**Tech Stack:** Go 1.22+, Gin, pgx, golang-migrate, golang.org/x/oauth2, PostgreSQL 16, Docker Compose

**Spec:** `docs/superpowers/specs/2026-06-11-backend-auth-profile-settings-design.md`

## Global Constraints

- No prompt/completion/request-response content storage
- Magic request/resend always HTTP 202 (no email enumeration)
- Unlink last auth method → `403 last_auth_method`
- Module path: `github.com/dxpglobal/ai-gateway/backend`
- Backend root: `lp/backend`
- Cookie: `dxp_session`, HttpOnly, Secure, SameSite=Lax; DB stores SHA-256 of token
- SMTP is production mail path; `MAIL_BACKEND=console` for tests only
- Supported default models: `deepseek-v4.1-flash`, `deepseek-v4-flash`, `mimo-v2.6-pro`, `mimo-v2.6-flash`, `glm-5.3-flash`, `glm-5.3`, `kimi-k3`, `minimax-m3`

---

### Task 1: Scaffold module, config, compose, migrations

**Files:**
- Create: `lp/backend/go.mod`
- Create: `lp/backend/cmd/server/main.go`
- Create: `lp/backend/internal/pkg/config/config.go`
- Create: `lp/backend/docker-compose.yml`
- Create: `lp/backend/migrations/0001_init.up.sql`, `0001_init.down.sql`
- Create: `lp/backend/Makefile`

**Interfaces:**
- Produces: `config.Load() (*config.Config, error)` with fields from spec §Config
- Produces: DB schema in spec §Data model

- [ ] `go mod init github.com/dxpglobal/ai-gateway/backend`
- [ ] Add deps: `gin`, `pgx/v5`, `golang-migrate`, `golang.org/x/oauth2`, `bcrypt`
- [ ] Write `0001_init.up.sql` matching spec tables (`citext` extension first)
- [ ] `docker-compose.yml` with Postgres 16, DB `dxp_ai_gateway`, user `dxp`
- [ ] `config.Load` from env with defaults (`HTTP_ADDR=:8080`, cookie `dxp_session`, TTL 30)
- [ ] `main.go` loads config, opens pool, runs migrate, starts gin with `/healthz`
- [ ] Run: `go build ./...` → PASS; `docker compose up -d` + migrate → PASS

- [ ] Commit: `feat(backend): scaffold gin module, config, migrations, compose`

### Task 2: Entities + repositories

**Files:**
- Create: `lp/backend/internal/entity/{user,session,settings,auth_method,magic_code}.go`
- Create: `lp/backend/internal/repository/repository.go` (interfaces)
- Create: `lp/backend/internal/repository/postgres/*.go`
- Test: `lp/backend/internal/repository/postgres/*_test.go`

**Interfaces:**
- Produces: `UserRepo`, `SessionRepo`, `SettingsRepo`, `AuthMethodRepo`, `MagicCodeRepo` with methods used by usecases (Create/GetByEmail, CreateSession/GetByHash/Revoke/RevokeAllForUser/ListActive, Get/UpsertSettings, UpsertAuthMethod/ListByUser/Delete, Insert/Consume/Invalidate)

- [ ] Entities as plain structs (uuid, time)
- [ ] Interfaces in `repository` package
- [ ] Postgres impls with pgx pool
- [ ] Unit-level tests with testcontainer or `DATABASE_URL` skip-if-unset
- [ ] Commit: `feat(backend): entities and postgres repositories`

### Task 3: Auth usecase — magic code + sessions

**Files:**
- Create: `lp/backend/internal/usecase/auth.go`
- Create: `lp/backend/internal/pkg/mail/{mailer.go,smtp.go,console.go,templates/*}`
- Test: `lp/backend/internal/usecase/auth_test.go` (fakes)

**Interfaces:**
- Produces: `AuthUsecase` with `RequestMagicCode`, `VerifyMagicCode`, `ResendMagicCode`, `Logout`, `CurrentUser`
- Consumes: repos + `mail.Mailer`

- [ ] Failing tests: request creates hashed code + sends mail; verify issues session; resend cooldown 60s; 5 attempts burn; register name only on create
- [ ] Implement `AuthUsecase` + bcrypt hash + 10 min TTL
- [ ] Templates `magic_code.{html,txt,subject}.tmpl`
- [ ] SMTP + console mailer
- [ ] Tests pass; commit `feat(backend): magic-code auth and mail templates`

### Task 4: OAuth (Google, GitHub) + logout

**Files:**
- Create: `lp/backend/internal/pkg/oauth/{google.go,github.go}`
- Modify: `lp/backend/internal/usecase/auth.go`
- Test: extend `auth_test.go`

**Interfaces:**
- Produces: `StartOAuth(provider)`, `HandleOAuthCallback(provider, code) → user + session`
- Produces: link/unlink connection methods on profile usecase (Task 5)

- [ ] Failing tests: new user create on callback; link existing by verified email; second provider links same user
- [ ] Implement oauth2 config + user info fetch
- [ ] Unlink last method → `ErrLastAuthMethod`
- [ ] Commit: `feat(backend): google and github oauth`

### Task 5: Profile + settings usecases

**Files:**
- Create: `lp/backend/internal/usecase/{profile.go,settings.go}`
- Test: `lp/backend/internal/usecase/{profile,settings}_test.go`

**Interfaces:**
- Produces: `ProfileUsecase` (Get/Update, Sessions list/revoke/revokeAll, Connections, Security flags)
- Produces: `SettingsUsecase` (Get/Patch)

- [ ] Failing tests for update name/email, revoke all including current, login_alerts, settings patch
- [ ] Implement usecases
- [ ] Commit: `feat(backend): profile and settings usecases`

### Task 6: Controllers + routes + middleware

**Files:**
- Create: `lp/backend/internal/controller/{auth,profile,settings,respond}.go`
- Create: `lp/backend/internal/middleware/session.go`
- Modify: `lp/backend/cmd/server/main.go`
- Test: `lp/backend/internal/controller/*_test.go`

**Interfaces:**
- Produces: HTTP API in spec §API surface under `/api/v1`

- [ ] Failing httptest: 401 without cookie; login sets cookie; logout clears
- [ ] Handlers + DTO validation
- [ ] `RequireSession` middleware
- [ ] CSRF header check on mutations
- [ ] Wire routes in main
- [ ] Commit: `feat(backend): http controllers and session middleware`

### Task 7: Frontend wiring

**Files:**
- Modify: `lp/frontend/src/components/app-shell.tsx` (Logout item)
- Modify: `lp/frontend/src/components/auth-form.tsx`
- Modify: `lp/frontend/src/app/profile/page.tsx`, `settings/page.tsx`
- Modify: `lp/frontend/next.config.ts` (rewrites `/api/auth/signin/*`)

- [ ] Logout button → POST logout → `/login`
- [ ] Magic form → request/verify/resend APIs
- [ ] Profile/Settings load + save
- [ ] Hide 2FA/password/phone; monogram avatar
- [ ] `npm run build` + `tsc --noEmit` from `lp/frontend` PASS
- [ ] Commit: `feat(frontend): wire auth, profile, settings, logout`

---

## Execution handoff

After Task 1–7: run `go test ./...` and frontend build. Manual smoke: magic register → profile → settings → logout → OAuth (if credentials present).
