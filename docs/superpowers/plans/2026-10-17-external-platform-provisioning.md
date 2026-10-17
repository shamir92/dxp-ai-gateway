# External Platform Provisioning Implementation Plan

**Goal:** Provision each newly created gateway user on TPAI and APAI, storing per-site passwords encrypted and retrying remote failures from PostgreSQL.

**Architecture:** Build two site-configured New API clients behind one request implementation. New local account creation will atomically persist the gateway user and per-site pending rows; a worker will claim rows with leases, provision users independently, and schedule retries.

**Tech Stack:** Go 1.26, Gin backend, PostgreSQL/pgx, golang-migrate, AES-256-GCM, `net/http`.

**Spec:** `docs/superpowers/specs/2026-10-17-external-platform-user-provisioning-design.md`

## Global Constraints

- Provision both sites on first local account creation, including first-time OAuth; do not provision on existing-user login or account linking.
- Use a separate random password for each site and encrypt it with AES-256-GCM using a base64 32-byte key from environment configuration.
- Never log credentials, request bodies, or response bodies.
- Retry network errors, HTTP 429, and 5xx; mark permanent 4xx failures for operator reconciliation.
- Access-token creation remains deferred.
- No tests or test commands in this execution, per the active session instruction.

## Review Focus

- Partial site configuration must not silently provision only one site.
- A stale worker lease must not overwrite the result of a newer claim.
- Auth database errors must not be mistaken for missing users and trigger duplicate provisioning.
- Remote redirects must not receive admin bearer credentials.
- Logs and persisted error summaries must not contain response bodies or secret values.

---

### Task 1: Configuration, entities, encryption, and New API client

**Files:**
- Modify: `lp/backend/internal/pkg/config/config.go`, `lp/backend/.env.example`
- Create: `lp/backend/internal/entity/external_site_account.go`
- Create: `lp/backend/internal/integration/credentials/cipher.go`
- Create: `lp/backend/internal/integration/credentials/factory.go`
- Create: `lp/backend/internal/integration/newapi/client.go`
- Create: `lp/backend/migrations/0003_external_site_accounts.up.sql`
- Create: `lp/backend/migrations/0003_external_site_accounts.down.sql`

**Produces:** Site configuration, credential factory/cipher, persistent account model, schema, and `newapi.Client.CreateUser(ctx, username, displayName string, password []byte) error`.

### Task 2: Atomic account creation and retry repository

**Files:**
- Modify: `lp/backend/internal/repository/repository.go`
- Create: `lp/backend/internal/repository/postgres/external_site_accounts.go`
- Modify: `lp/backend/internal/usecase/auth.go`, `auth_magic_code.go`, `auth_oauth.go`

**Consumes:** Task 1 entity and credential factory.

**Produces:** `NewAccountRepo.CreateNewAccount` transaction and lease-based external-account claim/result methods; auth hooks that enqueue accounts only for new users.

### Task 3: Worker and application wiring

**Files:**
- Create: `lp/backend/internal/worker/external_site_accounts.go`
- Modify: `lp/backend/cmd/server/main.go`
- Modify: `docs/superpowers/specs/2026-10-17-external-platform-user-provisioning-design.md`

**Consumes:** Task 1 site clients/cipher and Task 2 repository/auth integration.

**Produces:** Startup validation, worker lifecycle with graceful shutdown, retry and failure handling, and implementation notes for unresolved remote response/idempotency behavior.
