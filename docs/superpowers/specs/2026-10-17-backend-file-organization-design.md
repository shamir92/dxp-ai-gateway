# Backend Internal File Organization Design

**Date:** 2026-10-17  
**Status:** Draft for user review  
**Scope:** File organization under `lp/backend/internal`

## Goal

Make repository implementations, use case logic, and domain entities easier to locate by splitting the current large files into focused files. Preserve runtime behavior, package boundaries, dependency direction, and exported APIs.

## Current Structure

- `internal/repository/repository.go` contains repository interfaces.
- `internal/repository/postgres/postgres.go` contains implementations for users, auth methods, sessions, settings, magic codes, and login audit.
- `internal/usecase/auth.go` contains magic-code auth, OAuth, sessions, rate limiting, and token helpers.
- `internal/usecase/profile.go` contains profile and settings use cases.
- `internal/entity/entity.go` contains all domain entities and related helpers.

## Proposed Structure

### Repository

Keep repository interfaces in `internal/repository/repository.go`. Split the Postgres implementations into one file per repository, keeping the existing `postgres` package:

```text
internal/repository/postgres/
  users.go
  auth_methods.go
  sessions.go
  settings.go
  magic_codes.go
  login_audit.go
```

Each file owns its repository type, constructor, and SQL methods. Keep shared error behavior and query semantics unchanged.

### Use cases

Keep the existing `usecase` package and split auth methods by responsibility. Keep use case types, constructors, and package-level shared declarations in `auth.go`:

```text
internal/usecase/
  auth.go             # AuthUsecase, result/meta types, constructor, shared errors
  auth_magic_code.go  # request, resend, verify, code constants, and rate helpers
  auth_oauth.go       # OAuth completion, link, and unlink
  auth_session.go     # logout, session lookup/creation, login alerts, token helpers
  profile.go          # ProfileUsecase and profile, connection, and profile-session logic
  settings.go         # SettingsUsecase and settings patch logic
```

Methods remain on their current receiver types. This is a file organization change, not a division into new packages or use case instances.

### Entities

Replace `internal/entity/entity.go` with focused files in the existing `entity` package:

```text
internal/entity/
  user.go          # User and initials helpers
  auth_method.go   # Provider and AuthMethod
  session.go       # Session and Active
  settings.go      # UserSettings, region/model settings, and DefaultSettings
  magic_code.go    # MagicPurpose
  login_audit.go   # LoginAudit
```

Keep type names and fields unchanged so existing imports and consumers continue to work.

## Scope

The refactor covers only:

- `lp/backend/internal/repository/postgres`
- `lp/backend/internal/usecase`
- `lp/backend/internal/entity`

Leave migrations, controllers, middleware, and `internal/pkg` unchanged. Do not alter SQL schemas, API contracts, dependency wiring, or frontend code.

## Compatibility and Acceptance Criteria

- Existing package paths, exported type and function names, constructors, and interfaces remain unchanged.
- The Postgres SQL statements and repository behavior remain unchanged.
- Use case behavior and entity behavior remain unchanged.
- Each Postgres repository implementation and each entity definition has a clear, focused file.
- Auth, profile, and settings logic can be located by responsibility without changing receiver types or dependency injection.
- Remove `postgres.go` and `entity.go` after moving their contents. Keep `auth.go` for shared auth declarations and `profile.go` for profile logic.
- Existing worktree changes outside the stated scope remain untouched.

## Out of Scope

- Changing package boundaries or introducing domain subpackages.
- Changing repository interfaces, constructors, controller wiring, or dependency direction.
- Changing data models, migrations, SQL behavior, auth flows, or API behavior.
- Refactoring controllers, middleware, or `internal/pkg`.
- Adding functionality or changing frontend files.

## Existing Worktree Changes

At design time, the worktree contains user changes in frontend files, documentation related to the Pencil design implementation, and an untracked `test.yaml`. Those changes are outside this spec's scope and must remain untouched during implementation.
