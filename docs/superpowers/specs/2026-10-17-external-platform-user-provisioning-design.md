# External Platform User Provisioning

**Status:** Draft for user review  
**Date:** 2026-10-17

## Goal

When the gateway creates a new local user, provision that user on both TPAI and APAI. Store the credentials for each site so a later backend flow can use the plaintext password to create an access token, while keeping plaintext credentials out of the database, API responses, and logs.

## Current state

The Go backend creates local users in `AuthUsecase.VerifyMagicCode` and `AuthUsecase.CompleteOAuthLogin`. It stores users in PostgreSQL and has no external platform client, password field, access-token creation flow, or AWS Secrets Manager integration. Configuration is read from environment variables and local `.env` files.

## Approved decisions

- Provision both sites whenever a new gateway user is created, including first-time OAuth sign-up. Existing-user logins do not provision again.
- Use the gateway user's email as the remote `username` and their name as `display_name`, falling back to email when the name is empty.
- Use `role: 1`, matching the supplied create-user request.
- Generate a separate cryptographically random password for each site.
- Let local gateway sign-up complete if an external site is unavailable; retry provisioning in the background.
- Use a database-backed retry workflow in the existing backend, with no separate queue service.
- Defer access-token creation to a later phase. This phase stores credentials and provides an internal decryption capability for a future backend consumer.
- Read site URLs, admin bearer tokens, and the encryption key from environment variables. Use `.env` for local development. The production secret source is undecided until the deployment stack and cost are selected.

## Scope

### In scope

- A client for each configured site, posting to `POST /api/user/` with the site-specific admin bearer token.
- Durable per-site provisioning status and encrypted credentials in PostgreSQL.
- A background worker with independent retries for each site.
- Internal encryption and decryption of site passwords.
- Configuration, migration, startup wiring, and backend tests for the workflow.

### Out of scope

- Creating access tokens on either site.
- Frontend status pages or user-visible credential retrieval.
- Choosing a cloud provider or production secret-manager implementation.
- A public API to retry or inspect external provisioning.

## Architecture and data flow

1. An auth flow verifies a new account and identifies that the local user does not already exist.
2. The backend generates two independent passwords and encrypts each one.
3. In the same PostgreSQL transaction that creates the local user, it creates default settings, the initial auth method, and one pending external-account row for each site. Existing accounts and OAuth links to an existing user do not create new provisioning rows.
4. The normal gateway session is issued after the local transaction commits. The request does not wait for either external site.
5. A worker in the existing Go process claims due rows using a short database lease, then releases the transaction before making HTTP calls. It processes TPAI and APAI independently.
6. A successful response marks that site's row provisioned and saves the remote user ID if the response provides one. A failure updates only that site's retry state.
7. A future internal access-token use case may decrypt a site's password by local user ID and site. No HTTP endpoint exposes this operation.

Persisting provisioning rows in the local account transaction closes the gap where the local account exists but a process crash loses the request to provision the external accounts.

## Data model

Add an `external_site_accounts` table rather than password columns on `users`. A user can have two site-specific usernames, passwords, remote IDs, and provisioning states, so they should be represented as separate records.

Each row should include:

- `id`, `user_id` (foreign key with cascade delete), and `site` (`tpai` or `apai`)
- `username`
- `password_ciphertext` and its nonce, stored as binary data
- `encryption_key_version`
- `remote_user_id`, nullable until returned by the remote API
- `status` (`pending`, `retrying`, `provisioned`, or `failed`)
- `attempt_count`, `next_attempt_at`, and lease expiry/claim metadata
- A sanitized `last_error`, timestamps, and `provisioned_at`

Add a unique constraint on `(user_id, site)` and an index for due pending/retrying rows. Do not add credential fields to the public `entity.User` JSON representation.

## Encryption and secret handling

- Use AES-256-GCM with a 32-byte key supplied as base64 through `EXTERNAL_PASSWORD_ENCRYPTION_KEY`.
- Store the ciphertext and nonce together as binary values, plus a key version. Keep the key outside PostgreSQL.
- Bind ciphertext to the local user ID and site as authenticated additional data so ciphertext cannot be swapped between users or sites.
- Generate passwords with `crypto/rand`; decrypt only inside backend code when a future trusted operation needs them.
- Never log request or response bodies, admin tokens, plaintext passwords, or ciphertext. Persist only sanitized error summaries.
- `.env.example` may document the variable and format, but must not contain a usable key or admin token. Real `.env` files remain local secrets.
- Keep the app's configuration interface environment-based. Deployment may inject values from a managed secret store or another cost-appropriate source once the stack is selected.

## Configuration and external request

Use these environment variables:

- `TPAI_BASE_URL`, `TPAI_ADMIN_TOKEN`
- `APAI_BASE_URL`, `APAI_ADMIN_TOKEN`
- `EXTERNAL_PASSWORD_ENCRYPTION_KEY`

The configured base URL is fixed by deployment configuration, never supplied by an end user. The HTTP client uses HTTPS in deployed environments, a bounded timeout, and safe redirect handling so admin credentials cannot be forwarded to an unexpected host.

For each site, send the supplied request shape:

```json
{
  "username": "<gateway user email>",
  "display_name": "<gateway user name or email>",
  "password": "<site-specific generated password>",
  "role": 1
}
```

## Retry and failure handling

- Retry network errors, HTTP 429, and 5xx responses with exponential backoff and jitter; honor `Retry-After` when present.
- Treat validation failures and admin-token 401/403 responses as failed records requiring operator correction. Do not drop the encrypted credentials or retry them on every poll.
- Use a lease and atomic claim so multiple backend instances do not process the same due row concurrently. Expired leases become eligible again after a worker crash.
- Bound each HTTP request, cancel worker activity during graceful shutdown, and never hold a database transaction open during the network request.
- A timeout can happen after a remote account was created. The supplied contract does not establish idempotency, lookup-by-username, response fields, or duplicate-user behavior. Until verified, a duplicate conflict must be recorded for reconciliation, not assumed to be success.
- Define an internal operator procedure to correct configuration and requeue permanent failures before production. Do not add a public retry endpoint in this phase.

## Code areas

- `internal/pkg/config`: parse and validate the two site configurations and encryption key.
- `internal/integration` (new): typed site client and create-user request/response handling.
- `internal/entity` and `internal/repository`: provisioning records and repository interfaces.
- `internal/repository/postgres`: transactional account onboarding and external-account persistence/claim updates.
- `internal/usecase`: invoke onboarding for new magic-code and OAuth accounts, plus an internal credential decryption interface.
- `internal/worker` (new or equivalent package): poll, claim, call, retry, and update site records.
- `cmd/server/main.go`: wire configuration, clients, repository, and worker lifecycle.
- `migrations`: add and remove `external_site_accounts`.
- `.env.example`: document variable names and key format with blank placeholders only.

## Verification plan

- Verify new magic-code registration and first-time OAuth sign-up produce exactly one pending row per site; repeat login and OAuth linking to an existing user produce none.
- Verify credentials are distinct per site, encrypt/decrypt round-trips correctly, and tampering or using the wrong user/site fails authentication.
- Verify repository uniqueness, due-row claims, leases, and independent TPAI/APAI state transitions.
- Use fake HTTP servers to check paths, bearer headers, payload mapping, timeout behavior, redirect handling, and retry classification without contacting production sites.
- Verify logs and user-facing JSON never include tokens or credential material.
- Confirm the real sites' response schema, role semantics, idempotency/lookup behavior, and duplicate-user response before implementation.

## Open decisions

1. Which deployment secret source will provide the environment variables? Decide after the deployment stack and cost are evaluated.
2. What do the real create-user responses contain, and can a user be looked up or created idempotently by username? Confirm with the site API before implementation.
3. What operator mechanism should requeue permanent provisioning failures? This phase assumes a private runbook or internal operation, not a public API.
