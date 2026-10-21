// Command worker provisions external site accounts for newly registered users.
//
// It is a standalone module with its own copies of the small set of types it
// shares with the backend, so it builds and deploys independently of both the
// backend and the frontend. It owns no HTTP surface.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/worker/internal/config"
	"github.com/dxpglobal/ai-gateway/worker/internal/credentials"
	"github.com/dxpglobal/ai-gateway/worker/internal/entity"
	"github.com/dxpglobal/ai-gateway/worker/internal/newapi"
	"github.com/dxpglobal/ai-gateway/worker/internal/repository/postgres"
)

func main() {
	// Load .env before config (process env wins over file).
	if err := config.LoadDotEnv(".env"); err != nil && !os.IsNotExist(err) {
		log.Printf("dotenv: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	cipher, err := credentials.NewCipher(cfg.ExternalPasswordEncryptionKey)
	if err != nil {
		log.Fatalf("external credentials: %v", err)
	}

	tpaiClient, err := newapi.NewClient(cfg.TPAI.BaseURL, cfg.TPAI.AdminToken, 10*time.Second)
	if err != nil {
		log.Fatalf("TPAI client: %v", err)
	}
	apaiClient, err := newapi.NewClient(cfg.APAI.BaseURL, cfg.APAI.AdminToken, 10*time.Second)
	if err != nil {
		log.Fatalf("APAI client: %v", err)
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
		log.Fatalf("provisioner: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("provisioning worker: started")
	<-provisioner.Start(ctx)
	log.Printf("provisioning worker: stopped")
}
