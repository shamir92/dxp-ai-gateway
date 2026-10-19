# API Key Management & Subscription Usage — Implementation Plan

**Goal:** Let a signed-in gateway user manage their API keys and read their subscription usage
on both external platforms (TPAI, APAI) from the `/token-plan` page — by minting **one**
long-lived access token per `(user, site)`, storing it encrypted, and reusing it, with **no
repeated logins**.

**Architecture:** A per-site New API client gains user-bearer calls (`login` → `verify` →
`access_tokens`) plus token CRUD and subscription read. The minted `nap_…` token is encrypted
into the existing `external_site_accounts` row. A new usecase resolves the token (reusing the
stored one, minting lazily when absent/expired) and proxies the key operations. The
`/token-plan` page's existing static cards are wired to real data.

**Tech Stack:** Go 1.26, Gin, PostgreSQL/pgx, golang-migrate, AES-256-GCM, `net/http`;
Next 16 / React 19 / Tailwind v4.

**Spec:** `docs/superpowers/specs/2026-10-18-api-key-management-design.md`

## Global Constraints

- Mint **lazily** on first API-key/subscription request — never during signup.
- One access token per `(user, site)`; reuse it. No repeated `/api/user/login` calls.
- Token scopes: `api_key:read`, `api_key:write`, `wallet:read` (+ `api_key:reveal` for the
  one-time key reveal on create). Lifetime 30 days; re-mint on expiry.
- Credits factor is the single constant **`1428`**, applied to both `amount_total` and
  `amount_used`. Never inline it; `amount_total == 0` means unlimited.
- Never log tokens, passwords, request/response bodies, or the admin token. Persist only
  sanitized error summaries.
- Never hold a DB transaction open across an HTTP call.
- HTTPS origin only, bounded timeout, redirects not followed.
- Branch on the upstream `success` field, not on HTTP status.

## Review Focus

- Cipher AAD must separate the password ciphertext from the token ciphertext (distinct label),
  so the two can never be swapped for one another.
- Exactly **one** re-mint per request on upstream 401/403 — no loops.
- Upstream quirks handled: `auto_groups` can be `null`; the `DeletedAt` field is noise; a
  negative `remain_quota` with `unlimited_quota: true` must not render as a balance.
- No secret value ever appears in a log line or a `/api/v1/*` response body.
- The frontend receives a full key **only** as the response of a successful create.

---

### Task 1: New API client — user-bearer surface

**Files:**
- Modify: `lp/backend/internal/integration/newapi/client.go`
- Create: `lp/backend/internal/integration/newapi/login.go`, `token.go`, `subscription.go`

**Produces:**

- `Login(ctx, username string, password []byte) (SessionToken, error)` → `POST /api/user/login`
  (bearer = `data.access_token`; absent `access_token` = challenge → error).
- `Verify(ctx, bearer string, scope string, ctxPayload any, password []byte) (proof string, error)`
  → `POST /api/verify` (method `password`; proof = `data.proof_token`).
- `CreateAccessToken(ctx, bearer, proof, name string, scopes []string, expiresAt int64) (raw string, remoteID int, error)`
  → `POST /api/user/access_tokens` with header `X-Security-Proof`; `data.token` = `nap_…`.
- `ListTokens(ctx, bearer, page, size)` / `GetToken` / `RevealKey(ctx, bearer, id)` /
  `CreateToken(ctx, bearer, payload)` / `UpdateToken` / `DeleteToken` → `/api/token/*`.
- `SubscriptionSelf(ctx, bearer)` → `/api/subscription/self`.
- `Status(ctx)` → `/api/status` (`api_info[]`).

**Consumes:** the envelope `{success,message,data}`; existing HTTPS/no-redirect policy.

### Task 2: Token persistence — schema, entity, repository, AAD

**Files:**
- Create: `lp/backend/migrations/0004_external_site_access_token.up.sql` / `.down.sql`
- Modify: `lp/backend/internal/entity/external_site_account.go`
- Modify: `lp/backend/internal/repository/repository.go`
- Modify: `lp/backend/internal/repository/postgres/external_site_accounts.go`
- Modify: `lp/backend/internal/integration/credentials/cipher.go`

**Produces:** columns `access_token_ciphertext`, `access_token_nonce`,
`access_token_key_version`, `access_token_expires_at`, `access_token_remote_id`,
`access_token_status`; repository `GetAccessToken(userID, site)` /
`SaveAccessToken(userID, site, …)` / `InvalidateAccessToken(userID, site)`; a cipher entry point
that binds AAD to the `"access_token"` label + user ID + site.

**Consumes:** existing `credentials.Cipher` and the `external_site_accounts` row.

### Task 3: APIKeyUsecase — resolve token, then proxy

**Files:**
- Create: `lp/backend/internal/usecase/api_key.go`, `lp/backend/internal/usecase/subscription.go`

**Produces:**

- `ensureToken(ctx, userID, site)` → decrypt stored token; if missing/expired, run
  login → verify → create → persist; return the bearer.
- `ListKeys`, `CreateKey` (create then reveal **once**), `UpdateKey`, `DeleteKey`,
  `Subscription` — each resolving the bearer via `ensureToken`.
- `QuotaToCreditsFactor = 1428`, applied in the subscription view.

**Consumes:** Task 1 client, Task 2 repository, `credentials.Cipher`.

**Behaviour:** on upstream 401/403, invalidate the stored token and re-mint **once**; if that
also fails, return a typed error.

### Task 4: Controller, routes, wiring

**Files:**
- Create: `lp/backend/internal/controller/api_keys.go`
- Modify: `lp/backend/internal/controller/auth.go` (extend `writeErr` for upstream errors)
- Modify: `lp/backend/cmd/server/main.go`

**Produces:** session-scoped routes

```
GET    /api/v1/profile/api-keys?site=tpai|apai&p=&page_size=
POST   /api/v1/profile/api-keys
PUT    /api/v1/profile/api-keys/:site/:id
DELETE /api/v1/profile/api-keys/:site/:id
GET    /api/v1/profile/subscription?site=tpai|apai
```

`GET …/subscription` returns `credits_total` / `credits_used` (already ×1428) plus raw
`amount_total` / `amount_used`, `plan_id`, `status`, `end_time`, `next_reset_time`.

### Task 5: Frontend — make `/token-plan` live

**Files:**
- Modify: `lp/frontend/src/components/token-plan-ui.tsx` (accept props, drop hardcoded values)
- Modify: `lp/frontend/src/components/token-plan-dashboard.tsx` (fetch, pass down, loading/error)
- Modify: `lp/frontend/src/app/token-plan/page.tsx` (if a wrapper is needed)
- Modify: `lp/frontend/src/lib/api.ts` (only if a helper is missing)

**Produces:** `PlanUsageCard` ← subscription (`amount_used`/`amount_total` ×1428);
`PlanMetaStrip` ← plan/status/dates; `ApiKeyCard` ← `/api/token/` list + create (one-time
reveal) + delete; `PlanBenefitsCard` ← the 10 models.
All requests via `apiJSON` from `src/lib/api.ts` (cookie session + CSRF header).

**Platform is fixed per page, not switchable.** `TokenPlanDashboard` takes a required `site`
prop: `/token-plan` and `/consumption` pass `tpai`; the API plan page passes `apai`.
There is no runtime TPAI/APAI toggle.

### Task 6: Verification

**Backend:** `httptest` fakes for every upstream call (paths, headers, bearer, bodies,
timeouts, no-redirect); cipher round-trip incl. wrong-label failure; repository get/save;
usecase: mint-once-then-reuse, single re-mint on 401, no token in logs.
**Frontend:** `npx tsc --noEmit` and `npx eslint` on changed files.
**Manual:** one real TPAI user — confirm the stored token is reused across many requests with
no further `/api/user/login` calls.

---

## Already done (this change)

The 10-model catalog (§Plan benefits): `pricing-shared.tsx`, `app/page.tsx`,
`app/models/page.tsx`, `app/api-pricing/page.tsx`, plus the plan-benefit line in
`app/pricing/page.tsx`. `PlanBenefitsCard` population is part of Task 5.
