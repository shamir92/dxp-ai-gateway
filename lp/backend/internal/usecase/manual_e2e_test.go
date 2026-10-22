//go:build manual

// This file is NOT part of the normal build. It exercises the real production
// path against the live external platforms using locally provisioned accounts.
// Run it only deliberately:
//
//	cd lp/backend && set -a && . ./.env; set +a
//	go test ./internal/usecase/ -tags manual -run TestManualEndToEnd -v -count=1
//
// It never logs a password, an access token, or ciphertext.

package usecase_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/credentials"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/newapi"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository/postgres"
	"github.com/dxpglobal/ai-gateway/backend/internal/usecase"
)

func TestManualEndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
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

	sites := []struct {
		site       entity.ExternalSite
		baseURL    string
		adminToken string
	}{
		{entity.ExternalSiteTPAI, os.Getenv("TPAI_BASE_URL"), os.Getenv("TPAI_ADMIN_TOKEN")},
		{entity.ExternalSiteAPAI, os.Getenv("APAI_BASE_URL"), os.Getenv("APAI_ADMIN_TOKEN")},
	}

	for _, target := range sites {
		t.Run(string(target.site), func(t *testing.T) {
			if target.baseURL == "" || target.adminToken == "" {
				t.Skipf("%s env not configured", target.site)
			}

			var userID uuid.UUID
			if err := pool.QueryRow(ctx, `
				SELECT user_id FROM external_site_accounts
				WHERE site = $1 AND status = 'provisioned'
				ORDER BY created_at LIMIT 1`, target.site).Scan(&userID); err != nil {
				t.Fatalf("select provisioned account: %v", err)
			}
			t.Logf("target user: %s", userID)

			client, err := newapi.NewClient(target.baseURL, target.adminToken, os.Getenv("CF_CLIENT_ID"), os.Getenv("CF_CLIENT_SECRET"), 20*time.Second)
			if err != nil {
				t.Fatalf("newapi client: %v", err)
			}
			factory, err := credentials.NewAccountFactory(cipher, []entity.ExternalSite{
				entity.ExternalSiteTPAI, entity.ExternalSiteAPAI,
			})
			if err != nil {
				t.Fatalf("account factory: %v", err)
			}
			uc := usecase.NewAPIKeyUsecase(postgres.NewUserRepo(pool), accounts, accounts, cipher, factory, map[entity.ExternalSite]usecase.SiteClient{
				target.site: client,
			})

			// 1st call: mints via login → verify → access_tokens and stores it.
			first, err := uc.ListKeys(ctx, userID, target.site, 1, 20)
			if err != nil {
				t.Fatalf("ListKeys (mint): %v", err)
			}
			t.Logf("mint + list OK: total=%d returned=%d", first.Total, len(first.Items))

			// 2nd call: must reuse the stored token, with no further login.
			second, err := uc.ListKeys(ctx, userID, target.site, 1, 20)
			if err != nil {
				t.Fatalf("ListKeys (reuse): %v", err)
			}
			t.Logf("reuse + list OK: total=%d returned=%d", second.Total, len(second.Items))

			// Subscription read with the same stored token.
			usage, err := uc.Subscription(ctx, userID, target.site)
			if err != nil {
				t.Fatalf("Subscription: %v", err)
			}
			t.Logf("subscription OK: billing_preference=%q plans=%d", usage.BillingPreference, len(usage.Plans))
			for _, plan := range usage.Plans {
				t.Logf("  plan=%d status=%q amount_total=%d amount_used=%d credits_total=%v credits_used=%d unlimited=%t",
					plan.PlanID, plan.Status, plan.AmountTotal, plan.AmountUsed,
					plan.CreditsTotal, plan.CreditsUsed, plan.Unlimited)
			}

			// The stored credential must be encrypted at rest, never plaintext.
			var (
				status     string
				ciphertext []byte
				expiresAt  *time.Time
			)
			if err := pool.QueryRow(ctx, `
				SELECT access_token_status, access_token_ciphertext, access_token_expires_at
				FROM external_site_accounts
				WHERE user_id = $1 AND site = $2`, userID, target.site).
				Scan(&status, &ciphertext, &expiresAt); err != nil {
				t.Fatalf("read stored token: %v", err)
			}
			if status != string(entity.AccessTokenActive) {
				t.Errorf("access_token_status = %q, want %q", status, entity.AccessTokenActive)
			}
			if len(ciphertext) == 0 {
				t.Error("stored ciphertext is empty")
			}
			if bytes.HasPrefix(ciphertext, []byte("nap_")) {
				t.Error("stored credential is PLAINTEXT (starts with nap_) — must be encrypted")
			}
			if expiresAt == nil {
				t.Error("stored token has no expiry")
			} else {
				t.Logf("stored credential: status=%s ciphertext_bytes=%d expiry=%s (%.1f days out)",
					status, len(ciphertext), expiresAt.Format(time.RFC3339),
					time.Until(*expiresAt).Hours()/24)
			}
		})
	}
}
