# Backend M1 — Auth, Profile, Settings Design

**Date:** 2026-06-11  
**Status:** Approved design (pending user review of this file)  
**Scope:** First backend milestone for DXP AI Gateway (`lp/backend`)

## Goal

Ship a Go/Gin control-plane API that powers login/register (magic code + Google + GitHub), profile, settings, and logout for the existing Next.js frontend in `lp/frontend`. LLM proxy, plugins, token plan, and API pricing are out of scope.

## Stack

- Go 1.22+, Gin
- Clean architecture: `controller → usecase → repository → entity`
- PostgreSQL 16, `pgx`, `golang-migrate`
- `golang.org/x/oauth2` for Google/GitHub
- SMTP for magic-code email
- Docker Compose for local Postgres

## Architecture

```
lp/backend/
  cmd/server/main.go          # config, DI, migrate, run
  internal/
    entity/                   # pure domain types
    repository/               # interfaces + postgres impls
    usecase/                  # Auth, Profile, Settings
    controller/               # gin handlers + DTOs
    middleware/               # RequireSession
    pkg/config/
    pkg/mail/                 # Mailer + SMTP impl
    pkg/oauth/                # Google + GitHub
  migrations/
  docker-compose.yml
```

Rules:

- Dependency direction is strictly `controller → usecase → repository → entity`. No upward imports.
- Gin appears only in `controller` and `middleware`.
- Usecases depend on interfaces (repos, mailer), not concrete Postgres/SMTP.

## Data model (Postgres)

```sql
users (
  id            uuid PK,
  email         citext UNIQUE NOT NULL,
  name          text NOT NULL DEFAULT '',
  avatar_url    text NOT NULL DEFAULT '',
  created_at    timestamptz NOT NULL,
  updated_at    timestamptz NOT NULL
)

auth_methods (
  id               uuid PK,
  user_id          uuid FK → users ON DELETE CASCADE,
  provider         text NOT NULL,          -- 'email' | 'google' | 'github'
  provider_user_id text NOT NULL,          -- email address or OAuth subject
  email            citext NOT NULL,
  UNIQUE (provider, provider_user_id)
)

sessions (
  id            uuid PK,
  user_id       uuid FK → users ON DELETE CASCADE,
  user_agent    text NOT NULL DEFAULT '',
  ip            inet,
  created_at    timestamptz NOT NULL,
  last_seen_at  timestamptz NOT NULL,
  revoked_at    timestamptz               -- null = active
)

magic_codes (
  id          uuid PK,
  email       citext NOT NULL,
  code_hash   text NOT NULL,              -- hash of 6-digit code
  purpose     text NOT NULL,              -- 'login' | 'register'
  expires_at  timestamptz NOT NULL,       -- 10 minutes
  consumed_at timestamptz,
  attempts    int NOT NULL DEFAULT 0
)

user_settings (
  user_id            uuid PK FK → users ON DELETE CASCADE,
  default_region     text NOT NULL DEFAULT 'sgp',   -- 'sgp' | 'us'
  default_model      text NOT NULL DEFAULT 'mimo-v2.6-flash',
  notify_balance     boolean NOT NULL DEFAULT true,
  notify_usage       boolean NOT NULL DEFAULT true,
  notify_security    boolean NOT NULL DEFAULT true,
  notify_product     boolean NOT NULL DEFAULT false,
  stream_default     boolean NOT NULL DEFAULT true,
  store_request_logs boolean NOT NULL DEFAULT true,  -- metadata only (future)
  login_alerts       boolean NOT NULL DEFAULT true
)

login_audit (
  id          uuid PK,
  user_id     uuid FK → users ON DELETE CASCADE,
  session_id  uuid FK → sessions ON DELETE CASCADE,
  ip          inet,
  user_agent  text,
  created_at  timestamptz NOT NULL
)
```

### Session cookie

- Name: `dxp_session` (configurable)
- Value: raw opaque session UUID
- DB stores only SHA-256 hash of the token
- Attributes: `HttpOnly`, `Secure`, `SameSite=Lax`, `Path=/`
- TTL: configurable (default 30 days), enforced via `revoked_at` + max age on cookie

Middleware: hash cookie → load active session → load user → `c.Set("user", user)`.

## Identity model

- One `User` row per email identity.
- Multiple `auth_methods` rows (email magic / google / github).
- Unlink of the last remaining method is rejected (`403 last_auth_method`).
- Google/GitHub: if provider email matches an existing user and is verified by the provider, link to that user; otherwise create a new user.
- Magic code `purpose=register` with existing email still succeeds as login (does not overwrite `name`).
- Magic code `purpose=login` for unknown email creates the user on first successful verify (passwordless convenience). Full name is empty unless register provided it.

Every successful auth creates a `sessions` row and a `login_audit` row.

## API surface

Base: `/api/v1`. JSON errors: `{ "error": { "code": "...", "message": "..." } }`.

### Auth (public)

| Method | Path | Notes |
|--------|------|-------|
| POST | `/auth/magic/request` | `{ email, purpose: "login"\|"register", name? }` → 202 always (no enumeration). Register requires `name`. |
| POST | `/auth/magic/verify` | `{ email, code, purpose }` → create/link user → set cookie → `{ user }` |
| POST | `/auth/magic/resend` | `{ email, purpose }` — rate limited; invalidates previous code |
| GET | `/auth/oauth/google` | redirect to Google |
| GET | `/auth/oauth/google/callback` | upsert user + method → cookie → redirect `PUBLIC_WEB_URL/token-plan` |
| GET | `/auth/oauth/github` | redirect to GitHub |
| GET | `/auth/oauth/github/callback` | same as Google |
| POST | `/auth/logout` | revoke current session, clear cookie |
| GET | `/auth/me` | current user + settings; 401 if no session |

Frontend already calls `/api/auth/signin/{google,github}` — Next rewrite maps those to `/api/v1/auth/oauth/{google,github}`. Backend mounts OAuth only at `/api/v1/auth/oauth/*`.

### Profile (RequireSession)

| Method | Path | Notes |
|--------|------|-------|
| GET | `/profile` | user + connected providers + session count |
| PATCH | `/profile` | `{ name?, email? }` |
| POST | `/profile/avatar` | **Deferred** — monogram from name initials only in M1 |
| GET | `/profile/sessions` | device, ip, last_seen, `current` flag |
| DELETE | `/profile/sessions/:id` | revoke one |
| POST | `/profile/sessions/revoke-all` | revoke **all** sessions including current (UI "Sign out all") |
| GET | `/profile/connections` | `{ provider, connected, email }[]` |
| POST | `/profile/connections/:provider/link` | start OAuth link flow |
| DELETE | `/profile/connections/:provider` | unlink; 403 if last method |
| GET | `/profile/security` | `{ login_alerts }` |
| PATCH | `/profile/security` | `{ login_alerts: bool }` |

### Settings (RequireSession)

| Method | Path | Notes |
|--------|------|-------|
| GET | `/settings` | full settings row |
| PATCH | `/settings` | partial update of prefs |

### Middleware & CSRF

- `RequireSession` on `/profile/*`, `/settings`, `/auth/me`, `/auth/logout`
- Mutations require header `X-Requested-With: XMLHttpRequest` (or `X-CSRF-Protection: 1`) in addition to `SameSite=Lax` cookie

## Magic codes

- 6-digit numeric, hashed at rest with **bcrypt**
- TTL 10 minutes
- Max 5 verify attempts per code; then burn
- Resend cooldown 60s; resend invalidates previous code
- Rate limit: max 5 `request`+`resend` calls per email per hour (429 `rate_limited`)

## Email (SMTP)

`Mailer` interface used by usecases:

```go
type Mailer interface {
    SendMagicCode(ctx context.Context, to, code, purpose string) error
}
```

**SMTP implementation is the production path** (not console-only).

Config:

| Env | Purpose |
|-----|---------|
| `SMTP_HOST` | SMTP server host |
| `SMTP_PORT` | usually 587 or 465 |
| `SMTP_USERNAME` | auth user |
| `SMTP_PASSWORD` | auth password |
| `SMTP_FROM` | From header (e.g. `DXP AI Gateway <hello@fromdxp.com>`) |
| `SMTP_TLS` | `starttls` \| `tls` \| `none` (dev only) |

### Email templates (backend-owned)

Templates live in the backend, not in the frontend and not as ad-hoc strings in handlers.

```
lp/backend/internal/pkg/mail/templates/
  magic_code.html.tmpl
  magic_code.txt.tmpl
```

- Rendered with Go `html/template` + `text/template`
- Data: `Code`, `Purpose` (`login`/`register`), `ExpiresInMinutes`, `ProductName` (`DXP AI Gateway`), `SupportEmail` (`hello@fromdxp.com`)
- Subject from `magic_code.subject.tmpl` (plain text template)
- HTML body is the primary email; text/plain alternative always included
- Future templates (password reset, security alert) add files here; no template logic in controllers/usecases
- Optional `MAIL_BACKEND=console` still available for unit tests / local dry-run; default is `smtp`

Behavior:

- Body includes code, 10-minute expiry, and that the code is single-use
- Failures surface as `502 mail_send_failed` to the client (generic message); full error only in logs

## Domain rules summary

| Rule | Behavior |
|------|----------|
| Email enumeration | magic request/resend always 202 |
| First login create | unknown email + valid code → create user + settings |
| Register name | only applied when creating the user; never overwrites |
| OAuth link | match verified provider email to existing user |
| Last auth method | cannot unlink when it would leave zero methods |
| Privacy | never store prompts/completions/request-response content |
| Settings defaults | created in the same transaction as user create |

## Errors

| HTTP | code | When |
|------|------|------|
| 400 | `validation` | invalid body/params |
| 401 | `unauthenticated` | missing/invalid session |
| 403 | `last_auth_method` | unlink would leave no login method |
| 404 | `not_found` | missing session/resource |
| 409 | `conflict` | email/provider already linked to another user |
| 429 | `rate_limited` | magic request/resend |
| 502 | `mail_send_failed` | SMTP failure |
| 500 | `internal` | unexpected |

## Config (env)

| Var | Notes |
|-----|-------|
| `HTTP_ADDR` | default `:8080` |
| `DATABASE_URL` | Postgres DSN |
| `SESSION_COOKIE_NAME` | default `dxp_session` |
| `SESSION_TTL_DAYS` | default `30` |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` / `GOOGLE_REDIRECT_URL` | OAuth |
| `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET` / `GITHUB_REDIRECT_URL` | OAuth |
| `SMTP_*` | see Email section |
| `MAIL_BACKEND` | `smtp` (default) \| `console` |
| `PUBLIC_WEB_URL` | post-auth redirect base (e.g. `http://localhost:3000`) |

## Testing strategy (TDD)

1. **Usecase unit tests** — in-memory fakes for repos + mailer; table-driven
2. **Repository integration tests** — real Postgres (compose or testcontainers)
3. **Controller tests** — `httptest` + gin; cookie set/cleared; 401 without session
4. **Flows:** magic register, magic login, resend cooldown, expired/wrong code, Google/GitHub callback create+link, logout, revoke all, unlink last method, settings get/patch

## Frontend changes (same milestone)

| Change | File |
|--------|------|
| Logout in ACCOUNT sidebar → `POST /api/v1/auth/logout` → `/login` | `app-shell.tsx` |
| Wire magic form to request/verify/resend | `auth-form.tsx` |
| OAuth → `/api/auth/signin/{google,github}` + Next rewrite to backend | `auth-form.tsx`, Next config |
| Session gate on app shell via `/auth/me` | `app-shell.tsx` or layout |
| Profile + Settings real fetch/save | `profile/page.tsx`, `settings/page.tsx` |
| Hide/disable 2FA, password, phone rows; "Change photo" deferred (monogram) | `profile/page.tsx` |

## Out of scope

- Plugins, Token plan, API lifecycle (later renamed **API pricing**)
- 2FA, password auth, phone connect, Open Platform account
- Billing, API keys, usage metering, LLM proxy / new-api integration
- Real-time notification delivery (prefs are stored only)

## Definition of done

1. `go test ./...` green
2. Manual: magic-code register → profile edit → settings save → sidebar logout → Google/GitHub login
3. Unlink last auth method returns 403
4. Magic code email sends via SMTP (console backend acceptable in unit tests only)
5. No request/response content stored
