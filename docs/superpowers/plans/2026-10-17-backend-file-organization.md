# Backend Internal File Organization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Split backend repository implementations, use case logic, and entities into focused files without changing behavior or package APIs.

**Architecture:** Keep the existing `postgres`, `usecase`, and `entity` packages. Move declarations and methods into files grouped by repository or responsibility, preserving names, receiver types, constructors, SQL, and data definitions.

**Tech Stack:** Go, PostgreSQL, pgx.

**Spec:** [Backend Internal File Organization Design](../specs/2026-10-17-backend-file-organization-design.md)

## Global Constraints

- Existing package paths, exported type and function names, constructors, and interfaces remain unchanged.
- The Postgres SQL statements and repository behavior remain unchanged.
- Use case behavior and entity behavior remain unchanged.
- The refactor covers only `lp/backend/internal/repository/postgres`, `lp/backend/internal/usecase`, and `lp/backend/internal/entity`.
- Leave migrations, controllers, middleware, and `internal/pkg` unchanged.
- Do not alter SQL schemas, API contracts, dependency wiring, or frontend code.
- Existing worktree changes outside the stated scope remain untouched.

## Review Focus

- **Postgres query drift:** Query text, scan targets, and error mapping must move unchanged with each repository method.
- **Repository API drift:** Concrete repository types, constructors, and method signatures must retain their current names and signatures.
- **Use case drift:** Receiver types, constructors, package-level errors, and helper behavior must remain unchanged as methods move.
- **Entity drift:** JSON tags, fields, constants, helper methods, and default settings values must move unchanged.
- **Scope drift:** No files outside the three approved internal areas should be modified.

This plan covers file movement only. It adds no test files and does not include test execution.

---

### Task 1: Split the Postgres repositories

**Files:**
- Create: `lp/backend/internal/repository/postgres/users.go`
- Create: `lp/backend/internal/repository/postgres/auth_methods.go`
- Create: `lp/backend/internal/repository/postgres/sessions.go`
- Create: `lp/backend/internal/repository/postgres/settings.go`
- Create: `lp/backend/internal/repository/postgres/magic_codes.go`
- Create: `lp/backend/internal/repository/postgres/login_audit.go`
- Delete: `lp/backend/internal/repository/postgres/postgres.go`
- Leave unchanged: `lp/backend/internal/repository/repository.go`

**Interfaces:**
- Consumes: Existing repository interfaces and entity types.
- Produces: The same concrete types and constructors currently used by `cmd/server/main.go`: `UserRepo` / `NewUserRepo`, `AuthMethodRepo` / `NewAuthMethodRepo`, `SessionRepo` / `NewSessionRepo`, `SettingsRepo` / `NewSettingsRepo`, `MagicCodeRepo` / `NewMagicCodeRepo`, and `LoginAuditRepo` / `NewLoginAuditRepo`.

- [ ] Move `UserRepo`, its constructor and methods, and `scanUser` into `users.go`.
- [ ] Move `AuthMethodRepo` and its constructor and methods into `auth_methods.go`.
- [ ] Move `SessionRepo` and its constructor and methods into `sessions.go`.
- [ ] Move `SettingsRepo` and its constructor and methods into `settings.go`.
- [ ] Move `MagicCodeRepo` and its constructor and methods into `magic_codes.go`.
- [ ] Move `LoginAuditRepo` and its constructor and methods into `login_audit.go`.
- [ ] Remove `postgres.go` after all declarations have moved; retain the existing `postgres` package and add the imports each new file needs.
- [ ] Review the moved SQL strings, scan order, error handling, and public method signatures against the original file.
- [ ] Commit the repository file split as `refactor(backend): split postgres repositories`.

### Task 2: Split the entity declarations

**Files:**
- Create: `lp/backend/internal/entity/user.go`
- Create: `lp/backend/internal/entity/auth_method.go`
- Create: `lp/backend/internal/entity/session.go`
- Create: `lp/backend/internal/entity/settings.go`
- Create: `lp/backend/internal/entity/magic_code.go`
- Create: `lp/backend/internal/entity/login_audit.go`
- Delete: `lp/backend/internal/entity/entity.go`

**Interfaces:**
- Consumes: Existing entity declarations in `entity.go`.
- Produces: The same exported types, constants, fields, tags, and methods at the same `entity` package path.

- [ ] Move `User`, `Initials`, `splitWords`, and `upper` into `user.go`.
- [ ] Move `Provider`, its constants, and `AuthMethod` into `auth_method.go`.
- [ ] Move `Session` and `Active` into `session.go`.
- [ ] Move `UserSettings`, region constants, `AllowedModels`, and `DefaultSettings` into `settings.go`.
- [ ] Move `MagicPurpose` and its constants into `magic_code.go`.
- [ ] Move `LoginAudit` into `login_audit.go`.
- [ ] Remove `entity.go` after all declarations have moved; keep names, fields, JSON tags, values, and helper behavior unchanged.
- [ ] Commit the entity file split as `refactor(backend): split entity files`.

### Task 3: Split use case logic by responsibility

**Files:**
- Create: `lp/backend/internal/usecase/auth_magic_code.go`
- Create: `lp/backend/internal/usecase/auth_oauth.go`
- Create: `lp/backend/internal/usecase/auth_session.go`
- Create: `lp/backend/internal/usecase/settings.go`
- Modify: `lp/backend/internal/usecase/auth.go`
- Modify: `lp/backend/internal/usecase/profile.go`
- Leave unchanged: `lp/backend/internal/usecase/auth_test.go`

**Interfaces:**
- Consumes: Existing use case types, repository interfaces, entity types, and package-level errors.
- Produces: The existing `AuthUsecase`, `ProfileUsecase`, and `SettingsUsecase` types, constructors, methods, and receiver types.

- [ ] Keep `AuthUsecase`, `SessionMeta`, `AuthResult`, `NewAuthUsecase`, and shared error declarations in `auth.go`.
- [ ] Move magic code request, resend, and verify methods, their constants, `rateEntry`, `magicRate`, rate helpers, and `generateNumericCode` into `auth_magic_code.go`.
- [ ] Move `CompleteOAuthLogin`, `LinkOAuth`, and `UnlinkProvider` into `auth_oauth.go`.
- [ ] Move `Logout`, `SessionUser`, `issueSession`, login alert/device helpers, `generateToken`, and `hashToken` into `auth_session.go`.
- [ ] Keep `ProfileUsecase`, its view types, `connectionViews`, and all `ProfileUsecase` methods in `profile.go`.
- [ ] Move `SettingsPatch`, `SettingsUsecase`, its constructor, and its `Get` and `Patch` methods from `profile.go` into `settings.go`.
- [ ] Add per-file imports and remove imports no longer used by `auth.go` or `profile.go`; keep package-level unexported names unchanged.
- [ ] Review receiver types, method signatures, constants, and helper bodies against the original files.
- [ ] Commit the use case file split as `refactor(backend): split usecase files`.

### Task 4: Review the resulting file boundaries

**Files:**
- Review: all files created or modified in Tasks 1–3.
- Leave unchanged: all out-of-scope backend and frontend files.

- [ ] Confirm each Postgres repository implementation is in its corresponding file and `postgres.go` is removed.
- [ ] Confirm each entity group is in its corresponding file and `entity.go` is removed.
- [ ] Confirm auth responsibilities are in the planned files, `SettingsUsecase` is in `settings.go`, and existing receiver types remain the same.
- [ ] Confirm migrations, controllers, middleware, `internal/pkg`, frontend files, and unrelated worktree changes are untouched.
