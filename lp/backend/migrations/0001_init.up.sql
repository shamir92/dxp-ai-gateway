CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email       citext NOT NULL UNIQUE,
    name        text NOT NULL DEFAULT '',
    avatar_url  text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE auth_methods (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider         text NOT NULL CHECK (provider IN ('email', 'google', 'github')),
    provider_user_id text NOT NULL,
    email            citext NOT NULL,
    created_at       timestamptz NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_user_id)
);

CREATE INDEX auth_methods_user_id_idx ON auth_methods (user_id);

CREATE TABLE sessions (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash   text NOT NULL UNIQUE,
    user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_agent   text NOT NULL DEFAULT '',
    ip           inet,
    created_at   timestamptz NOT NULL DEFAULT now(),
    last_seen_at timestamptz NOT NULL DEFAULT now(),
    revoked_at   timestamptz
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);
CREATE INDEX sessions_token_hash_idx ON sessions (token_hash);

CREATE TABLE magic_codes (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email       citext NOT NULL,
    name        text NOT NULL DEFAULT '',
    code_hash   text NOT NULL,
    purpose     text NOT NULL CHECK (purpose IN ('login', 'register')),
    expires_at  timestamptz NOT NULL,
    consumed_at timestamptz,
    attempts    int NOT NULL DEFAULT 0,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX magic_codes_email_idx ON magic_codes (email);

CREATE TABLE user_settings (
    user_id            uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    default_region     text NOT NULL DEFAULT 'sgp',
    default_model      text NOT NULL DEFAULT 'mimo-v2.6-flash',
    notify_balance     boolean NOT NULL DEFAULT true,
    notify_usage       boolean NOT NULL DEFAULT true,
    notify_security    boolean NOT NULL DEFAULT true,
    notify_product     boolean NOT NULL DEFAULT false,
    stream_default     boolean NOT NULL DEFAULT true,
    store_request_logs boolean NOT NULL DEFAULT true,
    login_alerts       boolean NOT NULL DEFAULT true,
    updated_at         timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE login_audit (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id uuid NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    ip         inet,
    user_agent text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX login_audit_user_id_idx ON login_audit (user_id);
