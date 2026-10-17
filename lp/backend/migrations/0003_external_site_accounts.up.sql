CREATE TABLE external_site_accounts (
    id                     uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    site                   text NOT NULL CHECK (site IN ('tpai', 'apai')),
    username               citext NOT NULL,
    display_name           text NOT NULL,
    password_ciphertext    bytea NOT NULL,
    password_nonce         bytea NOT NULL,
    encryption_key_version text NOT NULL,
    remote_user_id         text,
    status                 text NOT NULL CHECK (status IN ('pending', 'retrying', 'provisioned', 'failed')),
    attempt_count          integer NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
    next_attempt_at        timestamptz NOT NULL DEFAULT now(),
    lease_token            uuid,
    lease_expires_at       timestamptz,
    last_error             text NOT NULL DEFAULT '',
    created_at             timestamptz NOT NULL DEFAULT now(),
    updated_at             timestamptz NOT NULL DEFAULT now(),
    provisioned_at         timestamptz,
    UNIQUE (user_id, site),
    UNIQUE (site, username)
);

CREATE INDEX external_site_accounts_due_idx
    ON external_site_accounts (next_attempt_at, created_at)
    WHERE status IN ('pending', 'retrying');
