# API Key Management (TPAI & APAI)

**Status:** Design — awaiting review
**Date:** 2026-10-18

## Goal

Let a signed-in DXP AI Gateway user list, create, edit, and delete their own API keys
(the New API "tokens") on both external sites (TPAI and APAI) from the gateway UI, and
read their **subscription usage** (quota total vs. used) from the same sites.

Each user's gateway identity maps to one provisioned account per site (already built).
We obtain **one long-lived access token per user per site, once**, store it encrypted,
and reuse it for every key-management call — **no repeated logins**.

## Current state

`feat: provision external platform accounts` (commit `c7c16bc`) creates the remote user on
both sites when a gateway user signs up, and stores the site password encrypted:

- `external_site_accounts` — one row per `(user, site)` with `password_ciphertext`,
  `password_nonce`, `encryption_key_version`, `status`, lease/retry metadata.
- `credentials.Cipher` (AES-256-GCM) encrypt/decrypt; `credentials.AccountFactory` generates
  a distinct random password per site.
- `newapi.Client.CreateUser` posts `POST /api/user/` with the **site admin token**.
- `worker.ExternalSiteAccountWorker` claims due rows, provisions, retries.

**Not built (verified, not assumed):**

- No access token is created or stored anywhere. `grep -riE "access_tokens?|security.?proof|/api/verify|user/login|X-Security" lp/backend` returns nothing.
- There are only three migrations (`0001`, `0002`, `0003`); `external_site_accounts` has no token column.
- `remote_user_id` exists but is always `NULL` — the worker calls `MarkProvisioned(..., "")`.
- `docs/superpowers/specs/2026-10-17-external-platform-user-provisioning-design.md` labels this work **"Deferred access-token flow"**, and its plan states *"Access-token creation remains deferred."*

So this spec **is** the deferred phase, now being built, plus the API-key endpoints on top.

## Approved decisions

- One access token per `(gateway user, site)`, created once, stored encrypted, reused.
- No repeated logins: after the token exists, key management never calls `/api/user/login` again.
- The token is scoped narrowly to what the gateway needs (API-key management + subscription read).
- Both TPAI and APAI, same implementation, per-site base URL + admin token.
- Backend endpoints **and** a frontend UI.
- The external target is upstream **new-api** v1.0.0-rc.42 (confirmed by the
  `x-new-api-version` response header on `GET /api/status`).
- **Mint lazily** on the user's first API-key or subscription request, not during signup.
- API keys are listed **masked**; the full key is revealed **once, on create** only.
- Access-token lifetime: **30 days**, re-minted on expiry.
- Credits factor is a single shared constant **`1428`**, applied to both `amount_total`
  and `amount_used`.

## External contract (new-api v1.0.0-rc.42)

Confirmed from the site's live responses and the upstream source at tag `v1.0.0-rc.42`.

### Response envelope

```json
{ "success": true, "message": "", "data": { } }
```

Errors carry `success:false` and a `message` (and sometimes `code`), often with HTTP 200.
**Always branch on `success`, never on HTTP status alone.**

### Phase A — mint the access token (requires a browser session, runs once)

A **personal access token cannot create another access token** — `service/access_token_scope.go`
maps `VerificationScopeAccessTokenGenerate → ""`, meaning PATs can never obtain that proof.
So minting must happen from a fresh login session, which is why we do it at provisioning/first use.

**A1. Login** — `POST /api/user/login?turnstile=`

```jsonc
// request
{ "username": "<gateway user email>", "password": "<decrypted site password>" }
// 200 response
{ "success": true, "message": "",
  "data": { "access_token": "<JWT>", "token_type": "Bearer",
            "access_expires_at": <unix>, "session": {...}, "user": {...} } }
```

- The `access_token` is the bearer for A2/A3.
- Also sets a refresh cookie (must be captured and forwarded within the mint sequence).
- `turnstile_check=false` on this deployment, so no turnstile token is needed.
- If `data.access_token` is absent, the login returned a *verification challenge* (2FA etc.),
  not a session — treat as a provisioning failure, not a success.
- Possibly requires the `new_api_has_session=1` cookie; mirror the browser headers
  (`Origin`, `Referer`, `Cookie`) as `CreateUser` already does.

**A2. Request a security proof** — `POST /api/verify`

```jsonc
// headers: Authorization: Bearer <A1 access_token>, Content-Type: application/json
// request
{ "method": "password", "scope": "access_token.generate",
  "context": { "scopes": ["api_key:read","api_key:write","wallet:read"], "expires_at": <unix> },
  "password": "<decrypted site password>" }
// 200 response
{ "success": true, "message": "",
  "data": { "proof_token": "<proof>", "expires_at": 0,
            "method": "password", "scope": "access_token.generate" } }
```

- `method` accepts only `"password"` and `"2fa"`; `context` is a raw JSON object.
- `AccessTokenGenerateContext` has **exactly** the keys `scopes` and `expires_at`.
- The proof is **single-use** (`SECURITY_PROOF_CONSUMED` on replay).

**A3. Create the access token** — `POST /api/user/access_tokens`

```jsonc
// headers: Authorization: Bearer <A1 access_token>, X-Security-Proof: <A2 proof_token>
// request — scopes and expires_at MUST byte-match the A2 context, or the server
// answers SECURITY_PROOF_CONTEXT_MISMATCH
{ "name": "dxp-gateway", "scopes": ["api_key:read","api_key:write","wallet:read"], "expires_at": <unix> }
// 200 response
{ "success": true, "data": { "token": "nap_<43 chars>", "item": { "id": 7, ... } } }
```

- Store `data.token` **encrypted**; keep `item.id` for later update/revoke.
- `name` 1–64 chars. `expires_at` must be `0` (never) or ≥ `now + 3600s`.
- Max **20** tokens per user (`AccessTokenMaxPerUser`); default lifetime **30 days**
  (`AccessTokenDefaultExpiryDays`).
- Because A3 requires a *browser session* and PATs cannot obtain `access_token.generate`,
  A1–A3 must run as one sequence with the same login session (capture the refresh cookie).

### Phase B — manage keys with the stored token (no login)

Send `Authorization: Bearer nap_<...>`. Required scopes: `api_key:read`, `api_key:write`
and `wallet:read` (subscription usage); `api_key:reveal` only if we surface full key values.

| Method | Path | Purpose | Response `data` |
|---|---|---|---|
| GET | `/api/token/?p=1&page_size=20` | list (paged) | `{items:[Token+auto_groups], page, page_size, total}` |
| GET | `/api/token/:id` | one key | masked `Token` |
| POST | `/api/token/` | create | **none** — `{success:true,message:""}` |
| POST | `/api/token/:id/key` | reveal full key | `{key:"sk-..."}` |
| PUT | `/api/token/` | update (body carries `id`) | masked `Token` |
| DELETE | `/api/token/:id` | delete | none |
| GET | `/api/subscription/self` | subscription usage | `{billing_preference, subscriptions:[…], all_subscriptions:[…]}` |

**Important behaviours**

- **Create returns no key.** The generated key is revealed only by `POST /api/token/:id/key`,
  which needs the `api_key:reveal` scope. To show a new key once, we call it immediately after create.
- List responses contain the **masked** key (`key[:4] + "**********" + key[-4:]`).
- To reveal the full key **exactly once, on create** (the agreed flow): `POST /api/token/`,
  then immediately `POST /api/token/:id/key`. The list never returns a full key.
- The response leaks a GORM artifact `"DeletedAt": null` — ignore it.
- `auto_groups` comes back as **`null`** when empty, not `[]` — tolerate both.
- When `unlimited_quota` is `true`, `remain_quota` may be **negative** (observed `-276313`);
  do not render it as a balance.
- Delete/update are owner-scoped server-side (`user_id` filter), so cross-user access is impossible.
- `GET /api/status` (public) returns `data.api_info[]` — the **API addresses** to display
  (`[{color, description, id, route, url}]`, e.g. `https://tp-sg.fromdxp.com`).
- `GET /api/user/access_tokens/scopes` returns the scope dictionary — usable to discover
  valid scope strings instead of hardcoding them.
- **Pagination:** page param is `p`; size param is `page_size` (falling back to `ps`, then
  `size`), default **10**, capped at **100**.
- **Token `status`:** `1` enabled, `2` disabled, `3` expired, `4` exhausted.
- `expired_time = -1` means the key never expires.
- Branch on the `success` field, not the HTTP status: `401`/`403` proof failures come back as
  `403` with a `code`, but many business errors are HTTP **200** with `success:false`.

### Token object fields (`/api/token/`)

`id, user_id, key, status, name, created_time, accessed_time, expired_time, remain_quota,
unlimited_quota, model_limits_enabled, model_limits, allow_ips (*string), used_quota, group,
cross_group_retry` (+ `auto_groups []string` on responses).
`expired_time = -1` means never. Create request body = the same fields (see the supplied cURL).

### Subscription usage (`/api/subscription/self`)

The **gateway backend** calls this on the user's behalf at the **per-site** new-api
(e.g. `https://tpai.shamirhusein.my.id`, and the same for APAI), authenticating with the
stored access token. The browser never talks to TPAI/APAI directly — it only calls our
`/api/v1/*`, which proxies. Response:

```jsonc
{ "success": true, "message": "",
  "data": {
    "billing_preference": "…",
    "subscriptions":     [ { "subscription": { /* UserSubscription */ } } ],  // active
    "all_subscriptions": [ { "subscription": { /* UserSubscription */ } } ]   // incl. expired
  } }
```

Each entry is **wrapped**: the payload is at `entry.subscription`. `UserSubscription` fields:

`id, user_id, plan_id, amount_total (int64), amount_used (int64), start_time, end_time, status,
source, last_reset_time, next_reset_time, upgrade_group, prev_user_group, downgrade_group,
allow_wallet_overflow, created_at, updated_at`

- `status`: `active` | `expired` | `cancelled`.
- `amount_total` **0 = unlimited**.
- `start_time` / `end_time` / `next_reset_time` are Unix seconds (`0` = none).

**Display conversion (our side):** `amount_total` and `amount_used` are New API quota units.
Multiply **both by `1428`** before display, so

```
credits_total = amount_total * 1428
credits_used  = amount_used  * 1428
```

Observed live (2026-10-18, TPAI): `amount_total: 2941177`, `amount_used: 368207` →
`credits_total = 4_200_000_756` (≈ the "4.2 Billion Credits" plan scale),
`credits_used = 525_799_596`.

Keep `1428` as a single named constant (e.g. `QuotaToCreditsFactor`) — do not inline it, and do
not apply it to `amount_total` when it is `0` (unlimited stays "unlimited", not `0`).

## Architecture and data flow

```
provisioning / first use
  password (decrypted) ──▶ A1 login ──▶ JWT + cookie
                             └▶ A2 /api/verify ──▶ proof
                                  └▶ A3 /api/user/access_tokens ──▶ nap_...
                                        └─▶ store encrypted in external_site_accounts

every later request
  session user ──▶ load account row ──▶ decrypt nap_ ──▶ Bearer nap_ ──▶ /api/token/ CRUD
                                                                       └▶ /api/subscription/self
```

1. The gateway session authorizes the request (existing `RequireSession` middleware).
2. `APIKeyUsecase` loads the user's `external_site_accounts` row for the requested site.
3. If the row has no access token (or it is expired/revoked), it runs the mint sequence once,
   persists the token, then continues.
4. It then calls `/api/token/` with the decrypted token.
5. On `401`/`403` from `/api/token/` (token revoked remotely or expired), it re-mints **once**
   and retries, then surfaces an error.

No plaintext token or password is ever logged or returned to the frontend.

## Self-healing missing external accounts

A user can be unusable in two distinct ways. Both are healed when the plan page is opened.

**Case 1 — no `external_site_accounts` row at all** (registered before provisioning existed, or
through a path that skipped it). `APIKeyUsecase.EnsureProvisioned` loads the user, rebuilds the
per-site rows via `credentials.AccountFactory`, and inserts them with
`ON CONFLICT (user_id, site) DO NOTHING`. Existing rows are never modified, so it is safe on
every page load.

**Case 2 — the row says `provisioned` but the remote account does not exist.** This is the
consequence of a real defect: **new-api returns business failures as HTTP 200 with
`success:false`**, and `CreateUser` only checked the HTTP status, so a failed creation (e.g.
`duplicate key value violates unique constraint "uni_users_username"`) was recorded as
`provisioned`. Observed live: our row said `provisioned` while an admin search for the username
returned **0** users, and login answered *"Username or password is incorrect, or user has been
banned"*.

Fix, in two parts:

- **`CreateUser` now checks the response envelope**, so a business failure is an error and is no
  longer persisted as a successful provisioning. (This supersedes the earlier note in
  `2026-10-17-…-design.md` that *"a 2xx response marks the row provisioned"*.)
- On a **non-retryable login rejection** during minting, the usecase calls
  `ExternalSiteAccountRepo.Requeue`, flipping the row back to `pending` and clearing
  `provisioned_at`, so the worker recreates the remote account. It then reports
  `409 external_account_missing` so the caller retries. The requeue is deliberately **not**
  triggered by transient (retryable) failures, since re-provisioning during an outage risks
  duplicate remote accounts. `Requeue` only touches rows that are currently `provisioned` and
  only while `attempt_count < 5`, so it cannot loop.

Called at:

1. **token-plan page load** — every `tp_*` request ensures the rows before querying, so merely
   opening the page heals both cases.
2. **profile save** (`PATCH /api/v1/profile`) — best effort and non-fatal, since the profile
   update has already succeeded.
3. **the subscription flow** — **not implemented yet**; it must ensure rows the same way when built.

Provisioning stays asynchronous (the worker owns it), so the first `tp_*` call after a heal still
returns `409`. The token-plan UI therefore retries briefly (6 attempts × 3 s) with a
"Setting up your account…" message instead of showing a hard error.

## Data model

Extend `external_site_accounts` (migration `0004`), keeping credentials on the per-site row:

- `access_token_ciphertext bytea` — `nap_...` encrypted with the same AES-256-GCM cipher.
- `access_token_nonce bytea`
- `access_token_key_version text`
- `access_token_expires_at timestamptz` — NULL/zero means never.
- `access_token_remote_id text` — `item.id`, for future update/revoke.
- `access_token_status text` — `none | active | invalid` (provisioning/serving state).

Reuse the existing cipher, but bind AAD to a distinct label (`"access_token"`) alongside
user ID + site, so a token ciphertext cannot be swapped with a password ciphertext.

## Security

- Encrypt the access token at rest with the same 32-byte key; never store plaintext.
- Decrypt only in-memory for the duration of a request; `clear()` after use (match the worker).
- Never log tokens, passwords, request/response bodies, or the admin token; persist only
  sanitized error summaries (reuse `safeErrorSummary`).
- Scope the token to `api_key:read` + `api_key:write` + `wallet:read` (+ `api_key:reveal` only if needed).
- Do not follow redirects; HTTPS origin only; bounded timeout (reuse the existing client policy).
- No plaintext token in any `/api/v1/*` response — mask to a hint (last 4 chars) if shown at all.
- The mint sequence must not hold a DB transaction open across HTTP calls.

## Configuration

No new required env vars if we reuse `EXTERNAL_PASSWORD_ENCRYPTION_KEY`, `TPAI_*`, `APAI_*`.
Optional knobs (defaults): access-token lifetime days, token name, and the scopes string.

## Retry and failure handling

- Mint failures retry with backoff and jitter, honouring `Retry-After` (reuse the worker's delay logic).
- A `401`/`403` on mint (bad password / admin misconfig) is **permanent** → mark and stop retrying.
- On `/api/token/` `401`/`403`, re-mint once; if it still fails, return a typed error to the caller.
- Missing/invalid stored token is not user-fatal: it degrades to "mint on demand".

## Code areas

- `internal/integration/newapi/client.go` — add `Login`, `Verify` (mint), `CreateAccessToken`,
  the token CRUD calls (`ListTokens`/`GetToken`/`CreateToken`/`UpdateToken`/`DeleteToken`/`RevealToken`),
  `SubscriptionSelf`, and `Status`. Split the user-bearer calls from the existing admin-bearer call.
- `internal/integration/credentials` — extend/reuse `Cipher` for token AAD.
- `internal/entity` — token fields on `ExternalSiteAccount`; a `Token` view type.
- `internal/repository` — get/set access-token methods.
- `internal/repository/postgres` — `0004` persistence.
- `internal/usecase` — new `APIKeyUsecase` (get-or-mint token, then CRUD).
- `internal/controller` — `APIKeyController` + routes; reuse `writeErr`.
- `cmd/server/main.go` — wiring.
- `migrations/0004_*` — add/remove token columns.
- `lp/frontend/src/` — wire the `/token-plan` page (`app/token-plan/page.tsx`,
  `components/token-plan-dashboard.tsx`, `components/token-plan-ui.tsx`) to the new endpoints via
  `src/lib/api.ts`; populate `PlanBenefitsCard` with the 10 models.

## API surface (proposed)

```
GET    /api/v1/profile/tp_api_keys?p=&page_size=
POST   /api/v1/profile/tp_api_keys
PUT    /api/v1/profile/tp_api_keys/:id
DELETE /api/v1/profile/tp_api_keys/:id
GET    /api/v1/profile/tp_subscription
```

All session-scoped to the current user.

**The paths use a neutral product slug (`tp_` = token plan) so the upstream platform is
never named in a URL.** The backing platform is bound per route by
`controller.WithSite(...)` — it is never read from the request — and it is also absent from
every response body (`site` was removed from the key-page, created-key and subscription
payloads). The API plan will mirror these as `ap_*`; those routes are not registered yet.

The subscription response exposes `credits_total` / `credits_used` (already ×1428) and the
raw `amount_total` / `amount_used`, plus `plan_id`, `status`, `end_time`, `next_reset_time`.

## Plan benefits and model catalog (frontend)

The plan **benefit** must name the models included. Canonical list (exact IDs, 10):

| # | Model ID | Vendor |
|---|---|---|
| 1 | `deepseek-v4.1-flash-fast` | DeepSeek |
| 2 | `deepseek-v4-flash-fast` | DeepSeek |
| 3 | `deepseek-v4.1-flash` | DeepSeek |
| 4 | `deepseek-v4-flash` | DeepSeek |
| 5 | `mimo-v2.6-pro` | MiMo |
| 6 | `mimo-v2.6-flash` | MiMo |
| 7 | `glm-5.3` | GLM |
| 8 | `glm-5.3-flash` | GLM |
| 9 | `kimi-k3` | Kimi |
| 10 | `minimax-m3` | MiniMax |

**Current state** — the catalog is duplicated across the frontend, and it is short by two:

- `components/pricing-shared.tsx` → `modelGroups` (3 groups, **8** models) — drives `ModelsSection`.
- `app/pricing/page.tsx` → `flagship: Plan["features"]`, spread into all four plans (Lite/Standard/Pro/Max);
  its first entry is `"V2.6 Flagship Model Access"`.
- `app/page.tsx` (8), `app/models/page.tsx` (8), `app/api-pricing/page.tsx` (6) — the same list again.

**Change (decided)**

- Add a benefit line naming the models to the `flagship` `features` array in
  `app/pricing/page.tsx`, so all four plan cards (Lite / Standard / Pro / Max) show them.
- Populate `PlanBenefitsCard` (`components/token-plan-ui.tsx`) on the `/token-plan` page with the
  same 10 models — it currently says only `"deepseek, mimo, glm, kimi, minimax"`.
- Add `deepseek-v4.1-flash-fast` and `deepseek-v4-flash-fast` to **all four** catalog copies
  (`pricing-shared.tsx`, `app/page.tsx`, `app/models/page.tsx`, `app/api-pricing/page.tsx`) so
  each lists all 10.
- Keep the four copies in sync by hand for now; consolidating into one shared constant is
  deferred (not in this change).

## Verification plan

- Mint once per `(user, site)`; verify exactly one access token is stored, encrypted, and reused
  across many key operations with **no** further `/api/user/login` calls.
- Encrypt/decrypt round-trip; wrong user/site/label fails authentication.
- `api_key:read` list, `api_key:write` create/update/delete, `api_key:reveal` full key — each
  rejected without the right scope.
- Create-then-reveal returns the key exactly once; list never leaks a full key.
- Subscription: `amount_total`/`amount_used` are unwrapped from `entry.subscription` and scaled
  by exactly ×1428; `amount_total == 0` renders as unlimited, not `0`.
- A revoked/expired stored token triggers exactly one re-mint and retry, then a clean error.
- Fake HTTP servers assert paths, headers, bearer, payload mapping, timeouts, no-redirect.
- Logs and user-facing JSON never contain tokens, passwords, or ciphertext.

## Decisions

Resolved (2026-10-18):

1. **Mint timing** — lazily, on the user's first API-key or subscription request (not at signup).
2. **Key visibility** — masked in lists; the full key is revealed **once, on create**.
3. **Token lifetime** — 30 days, re-minted on expiry.
4. **Credits factor** — the fixed constant `1428`, applied to both `amount_total` and `amount_used`.
5. **`/api/subscription/self` and PAT scopes** — the whole `/api` group carries
   `middleware.AccessTokenAudit()`, and subscription routes use `UserAuth()`; we request
   `wallet:read` up front and will confirm scope enforcement against the live site before
   relying on it.

Resolved (2026-10-18) — placement:

6. **All of this lands on the `/token-plan` page** — `app/token-plan/page.tsx` →
   `components/token-plan-dashboard.tsx` → `components/token-plan-ui.tsx`. Its already-built
   shells become live:
   - `PlanUsageCard` ← `/api/subscription/self` (`amount_used` / `amount_total` × 1428)
   - `PlanMetaStrip` ← subscription `plan_id`, `status`, `start_time`, `end_time`
   - `ApiKeyCard` ← `/api/token/` (list / create / reveal-once / delete)
   - `PlanBenefitsCard` ← the 10-model list
   - the model filter and usage bar chart stay static for now (no usage-analytics endpoint
     is in scope)

7. **The platform is fixed per page — there is no TPAI/APAI switcher.**
   The **token plan** is served by **TPAI**; the **API plan** is served by **APAI**.
   `TokenPlanDashboard` therefore takes a required `site` prop: `/token-plan` and
   `/consumption` pass `tpai` (route `app/token-plan`); the API plan page is
   `app/api-plan` (renamed from `api-lifecycle`) and passes `apai`.

Everything else in this spec is confirmed against **captured live responses** and the upstream
`new-api` source at tag `v1.0.0-rc.42`.
