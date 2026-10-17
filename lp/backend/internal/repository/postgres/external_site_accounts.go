package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
)

type ExternalSiteAccountRepo struct {
	pool *pgxpool.Pool
}

func NewExternalSiteAccountRepo(pool *pgxpool.Pool) *ExternalSiteAccountRepo {
	return &ExternalSiteAccountRepo{pool: pool}
}

func (r *ExternalSiteAccountRepo) CreateNewAccount(ctx context.Context, user *entity.User, method *entity.AuthMethod, accounts []entity.ExternalSiteAccount) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if method == nil {
		return fmt.Errorf("new account auth method is required")
	}
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now
	method.UserID = user.ID
	if method.ID == uuid.Nil {
		method.ID = uuid.New()
	}
	method.CreatedAt = now

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin new account transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		INSERT INTO users (id, email, name, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Email, user.Name, user.AvatarURL, user.CreatedAt, user.UpdatedAt); err != nil {
		return fmt.Errorf("insert new user: %w", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO user_settings (user_id) VALUES ($1)`, user.ID); err != nil {
		return fmt.Errorf("insert default user settings: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO auth_methods (id, user_id, provider, provider_user_id, email, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		method.ID, user.ID, method.Provider, method.ProviderUserID, method.Email, method.CreatedAt); err != nil {
		return fmt.Errorf("insert initial auth method: %w", err)
	}
	for i := range accounts {
		account := &accounts[i]
		if account.ID == uuid.Nil {
			account.ID = uuid.New()
		}
		if account.UserID != user.ID {
			return fmt.Errorf("external account user does not match new user")
		}
		account.CreatedAt = now
		account.UpdatedAt = now
		if account.Status == "" {
			account.Status = entity.ExternalAccountPending
		}
		if account.NextAttemptAt.IsZero() {
			account.NextAttemptAt = now
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO external_site_accounts (
				id, user_id, site, username, display_name, password_ciphertext, password_nonce,
				encryption_key_version, status, attempt_count, next_attempt_at, last_error, created_at, updated_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
			account.ID, account.UserID, account.Site, account.Username, account.DisplayName,
			account.PasswordCiphertext, account.PasswordNonce, account.EncryptionKeyVersion,
			account.Status, account.AttemptCount, account.NextAttemptAt, account.LastError,
			account.CreatedAt, account.UpdatedAt); err != nil {
			return fmt.Errorf("insert external site account: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit new account transaction: %w", err)
	}
	return nil
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

var _ repository.NewAccountRepo = (*ExternalSiteAccountRepo)(nil)
var _ repository.ExternalSiteAccountRepo = (*ExternalSiteAccountRepo)(nil)
