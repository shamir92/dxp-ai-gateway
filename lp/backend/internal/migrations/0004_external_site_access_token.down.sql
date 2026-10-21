ALTER TABLE external_site_accounts
    DROP CONSTRAINT IF EXISTS external_site_accounts_access_token_status_check;

ALTER TABLE external_site_accounts
    DROP COLUMN IF EXISTS access_token_ciphertext,
    DROP COLUMN IF EXISTS access_token_nonce,
    DROP COLUMN IF EXISTS access_token_key_version,
    DROP COLUMN IF EXISTS access_token_expires_at,
    DROP COLUMN IF EXISTS access_token_remote_id,
    DROP COLUMN IF EXISTS access_token_status;
