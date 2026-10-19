//go:build manual

// Not part of the normal build. Proves the missing-account backfill against the
// real database without touching any external platform:
//
//	cd lp/backend && set -a && . ./.env; set +a
//	go test ./internal/usecase/ -tags manual -run TestManualBackfill -v -count=1
//
// It creates a throwaway user, and dates the backfilled rows one hour into the
// future so the provisioning worker can never claim them (ClaimDue only takes
// rows with next_attempt_at <= now). The user is deleted at the end, which
// cascades the rows away.

package usecase_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/credentials"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository/postgres"
)

func TestManualBackfill(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")
	keyText := os.Getenv("EXTERNAL_PASSWORD_ENCRYPTION_KEY")
	if dbURL == "" || keyText == "" {
		t.Fatal("missing env — source lp/backend/.env first")
	}
	key, err := base64.StdEncoding.DecodeString(keyText)
	if err != nil {
		t.Fatalf("decode encryption key: %v", err)
	}
	cipher, err := credentials.NewCipher(key)
	if err != nil {
		t.Fatalf("cipher: %v", err)
	}
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("connect db: %v", err)
	}
	defer pool.Close()
	accounts := postgres.NewExternalSiteAccountRepo(pool)

	// A user that has no external accounts, standing in for the edge case.
	email := fmt.Sprintf("heal-test-%d@example.invalid", time.Now().UnixNano())
	var userID uuid.UUID
	if err := pool.QueryRow(ctx, `
		INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id`,
		email, "Heal Test").Scan(&userID); err != nil {
		t.Fatalf("insert throwaway user: %v", err)
	}
	defer func() {
		if _, err := pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID); err != nil {
			t.Errorf("cleanup throwaway user: %v", err)
		}
	}()
	t.Logf("throwaway user: %s (no external accounts)", userID)

	// Date the rows an hour out so the worker never picks them up.
	notDue := time.Now().UTC().Add(time.Hour)
	build := func(password string) []entity.ExternalSiteAccount {
		rows := make([]entity.ExternalSiteAccount, 0, 2)
		for _, site := range []entity.ExternalSite{entity.ExternalSiteTPAI, entity.ExternalSiteAPAI} {
			ct, nonce, kv, err := cipher.Encrypt(userID, site, []byte(password))
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			rows = append(rows, entity.ExternalSiteAccount{
				UserID:               userID,
				Site:                 site,
				Username:             email,
				DisplayName:          "Heal Test",
				PasswordCiphertext:   ct,
				PasswordNonce:        nonce,
				EncryptionKeyVersion: kv,
				Status:               entity.ExternalAccountPending,
				NextAttemptAt:        notDue,
			})
		}
		return rows
	}

	countRows := func() int {
		var n int
		if err := pool.QueryRow(ctx,
			`SELECT count(*) FROM external_site_accounts WHERE user_id = $1`, userID).Scan(&n); err != nil {
			t.Fatalf("count rows: %v", err)
		}
		return n
	}

	if before := countRows(); before != 0 {
		t.Fatalf("precondition failed: throwaway user already has %d rows", before)
	}

	// 1. First heal creates both rows.
	first := build("first-password")
	if err := accounts.EnsurePending(ctx, first); err != nil {
		t.Fatalf("EnsurePending (first): %v", err)
	}
	if got := countRows(); got != 2 {
		t.Fatalf("rows after first heal = %d, want 2 (one per site)", got)
	}
	t.Logf("first heal created %d rows (tpai + apai)", countRows())

	// 2. A second heal must not duplicate or overwrite anything.
	var originalCipher []byte
	if err := pool.QueryRow(ctx, `
		SELECT password_ciphertext FROM external_site_accounts
		WHERE user_id = $1 AND site = 'tpai'`, userID).Scan(&originalCipher); err != nil {
		t.Fatalf("read original ciphertext: %v", err)
	}

	if err := accounts.EnsurePending(ctx, build("second-password")); err != nil {
		t.Fatalf("EnsurePending (second): %v", err)
	}
	if got := countRows(); got != 2 {
		t.Fatalf("rows after second heal = %d, want 2 (no duplicates)", got)
	}

	var afterCipher []byte
	if err := pool.QueryRow(ctx, `
		SELECT password_ciphertext FROM external_site_accounts
		WHERE user_id = $1 AND site = 'tpai'`, userID).Scan(&afterCipher); err != nil {
		t.Fatalf("read ciphertext after second heal: %v", err)
	}
	if !bytes.Equal(originalCipher, afterCipher) {
		t.Error("second heal OVERWROTE the existing row — ON CONFLICT DO NOTHING is not holding")
	}
	t.Logf("second heal was a no-op: still %d rows, original row untouched", countRows())
}
