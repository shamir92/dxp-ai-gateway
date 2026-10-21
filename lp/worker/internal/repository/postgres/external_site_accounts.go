package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/worker/internal/entity"
	"github.com/dxpglobal/ai-gateway/worker/internal/repository"
)

// ExternalSiteAccountRepo is the worker's view of the external_site_accounts
// queue. It only carries the operations provisioning needs: claim a due row by
// lease, then finish it. Account creation and token storage stay in the backend.
type ExternalSiteAccountRepo struct {
	pool *pgxpool.Pool
}

func NewExternalSiteAccountRepo(pool *pgxpool.Pool) *ExternalSiteAccountRepo {
	return &ExternalSiteAccountRepo{pool: pool}
}

func (r *ExternalSiteAccountRepo) ClaimDue(ctx context.Context, now, leaseExpiresAt time.Time, limit int) ([]entity.ExternalSiteAccount, error) {
	rows, err := r.pool.Query(ctx, `
		WITH due AS (
			SELECT id
			FROM external_site_accounts
			WHERE status IN ('pending', 'retrying')
			  AND next_attempt_at <= $1
			  AND (lease_expires_at IS NULL OR lease_expires_at <= $1)
			ORDER BY next_attempt_at, created_at
			LIMIT $3
			FOR UPDATE SKIP LOCKED
		)
		UPDATE external_site_accounts AS account
		SET status = 'retrying', lease_token = gen_random_uuid(), lease_expires_at = $2, updated_at = $1
		FROM due
		WHERE account.id = due.id
		RETURNING account.id, account.user_id, account.site, account.username, account.display_name,
		          account.password_ciphertext, account.password_nonce, account.encryption_key_version,
		          account.remote_user_id, account.status, account.attempt_count, account.next_attempt_at,
		          account.lease_token, account.lease_expires_at, account.last_error, account.created_at,
		          account.updated_at, account.provisioned_at`, now, leaseExpiresAt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]entity.ExternalSiteAccount, 0)
	for rows.Next() {
		account, err := scanExternalSiteAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, rows.Err()
}

func (r *ExternalSiteAccountRepo) MarkProvisioned(ctx context.Context, id, leaseToken uuid.UUID, remoteUserID string, at time.Time) error {
	return r.finishClaim(ctx, `
		UPDATE external_site_accounts
		SET status = 'provisioned', remote_user_id = NULLIF($3, ''), attempt_count = attempt_count + 1,
		    last_error = '', lease_token = NULL, lease_expires_at = NULL, provisioned_at = $4, updated_at = $4
		WHERE id = $1 AND lease_token = $2`, id, leaseToken, remoteUserID, at)
}

func (r *ExternalSiteAccountRepo) MarkRetry(ctx context.Context, id, leaseToken uuid.UUID, nextAttemptAt time.Time, summary string, at time.Time) error {
	return r.finishClaim(ctx, `
		UPDATE external_site_accounts
		SET status = 'retrying', attempt_count = attempt_count + 1, next_attempt_at = $3,
		    last_error = $4, lease_token = NULL, lease_expires_at = NULL, updated_at = $5
		WHERE id = $1 AND lease_token = $2`, id, leaseToken, nextAttemptAt, summary, at)
}

func (r *ExternalSiteAccountRepo) MarkFailed(ctx context.Context, id, leaseToken uuid.UUID, summary string, at time.Time) error {
	return r.finishClaim(ctx, `
		UPDATE external_site_accounts
		SET status = 'failed', attempt_count = attempt_count + 1, last_error = $3,
		    lease_token = NULL, lease_expires_at = NULL, updated_at = $4
		WHERE id = $1 AND lease_token = $2`, id, leaseToken, summary, at)
}

func (r *ExternalSiteAccountRepo) finishClaim(ctx context.Context, query string, args ...any) error {
	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrLeaseLost
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanExternalSiteAccount(row rowScanner) (entity.ExternalSiteAccount, error) {
	var account entity.ExternalSiteAccount
	var remoteUserID pgtype.Text
	var leaseExpiresAt pgtype.Timestamptz
	var provisionedAt pgtype.Timestamptz
	err := row.Scan(
		&account.ID, &account.UserID, &account.Site, &account.Username, &account.DisplayName,
		&account.PasswordCiphertext, &account.PasswordNonce, &account.EncryptionKeyVersion,
		&remoteUserID, &account.Status, &account.AttemptCount, &account.NextAttemptAt,
		&account.LeaseToken, &leaseExpiresAt, &account.LastError, &account.CreatedAt,
		&account.UpdatedAt, &provisionedAt,
	)
	if err != nil {
		return entity.ExternalSiteAccount{}, err
	}
	if remoteUserID.Valid {
		account.RemoteUserID = &remoteUserID.String
	}
	if leaseExpiresAt.Valid {
		account.LeaseExpiresAt = &leaseExpiresAt.Time
	}
	if provisionedAt.Valid {
		account.ProvisionedAt = &provisionedAt.Time
	}
	return account, nil
}

var _ repository.ExternalSiteAccountRepo = (*ExternalSiteAccountRepo)(nil)
