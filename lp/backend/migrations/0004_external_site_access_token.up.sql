ALTER TABLE external_site_accounts
    ADD COLUMN access_token_ciphertext  bytea,
    ADD COLUMN access_token_nonce       bytea,
    ADD COLUMN access_token_key_version text,
    ADD COLUMN access_token_expires_at  timestamptz,
    ADD COLUMN access_token_remote_id   text,
    ADD COLUMN access_token_status      text NOT NULL DEFAULT 'none';

ALTER TABLE external_site_accounts
    ADD CONSTRAINT external_site_accounts_access_token_status_check
    CHECK (access_token_status IN ('none', 'active', 'invalid'));
