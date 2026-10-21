// Command worker provisions external site accounts for newly registered users.
//
// It is a standalone module with its own copies of the small set of types it
// shares with the backend, so it builds and deploys independently of both the
// backend and the frontend. It owns no HTTP surface.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/worker/internal/config"
	"github.com/dxpglobal/ai-gateway/worker/internal/credentials"
	"github.com/dxpglobal/ai-gateway/worker/internal/entity"
	"github.com/dxpglobal/ai-gateway/worker/internal/newapi"
	"github.com/dxpglobal/ai-gateway/worker/internal/observability"
	"github.com/dxpglobal/ai-gateway/worker/internal/repository/postgres"
)

func main() {
	// Load .env before config (process env wins over file).
	if err := config.LoadDotEnv(".env"); err != nil && !os.IsNotExist(err) {
		slog.Warn("dotenv load failed", slog.Any("error", err))
	}

	cfg, err := config.Load()
	if err != nil {
		fatal("config", err)
	}

	// Sentry is optional: without SENTRY_DSN this is a no-op.
	observability.Init()
	defer observability.Flush()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		fatal("db", err)
	}
	defer pool.Close()

	cipher, err := credentials.NewCipher(cfg.ExternalPasswordEncryptionKey)
	if err != nil {
		fatal("external credentials", err)
	}

	tpaiClient, err := newapi.NewClient(cfg.TPAI.BaseURL, cfg.TPAI.AdminToken, 10*time.Second)
	if err != nil {
		fatal("TPAI client", err)
	}
	apaiClient, err := newapi.NewClient(cfg.APAI.BaseURL, cfg.APAI.AdminToken, 10*time.Second)
	if err != nil {
		fatal("APAI client", err)
	}

	provisioner, err := NewProvisioner(
		postgres.NewExternalSiteAccountRepo(pool),
		cipher,
		map[entity.ExternalSite]UserCreator{
			entity.ExternalSiteTPAI: tpaiClient,
			entity.ExternalSiteAPAI: apaiClient,
		},
	)
	if err != nil {
		fatal("provisioner", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("provisioning worker started")
	<-provisioner.Start(ctx)
	slog.Info("provisioning worker stopped")
}

// fatal reports a startup failure and exits non-zero, so a supervisor sees the
// failure instead of a running-but-idle worker.
func fatal(what string, err error) {
	slog.Error("startup failed", slog.String("stage", what), slog.Any("error", err))
	os.Exit(1)
}
