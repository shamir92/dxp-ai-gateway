package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/backend/internal/controller"
	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/credentials"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/newapi"
	"github.com/dxpglobal/ai-gateway/backend/internal/middleware"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/config"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/geoip"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/mail"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/oauth"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository/postgres"
	"github.com/dxpglobal/ai-gateway/backend/internal/usecase"
	"github.com/dxpglobal/ai-gateway/backend/internal/worker"
)

func main() {
	// Load lp/backend/.env before config (process env wins over file).
	if err := config.LoadDotEnv(".env"); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("dotenv: %v", err)
		} else if err := config.LoadDotEnv("lp/backend/.env"); err != nil && !os.IsNotExist(err) {
			log.Printf("dotenv: %v", err)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	tmplDir := os.Getenv("MAIL_TEMPLATES_DIR")
	if tmplDir == "" {
		tmplDir = "internal/pkg/mail/templates"
	}
	tmpl, err := mail.LoadTemplates(tmplDir)
	if err != nil {
		// fall back relative to executable cwd
		alt := filepath.Join("lp", "backend", "internal", "pkg", "mail", "templates")
		tmpl, err = mail.LoadTemplates(alt)
		if err != nil {
			log.Fatalf("mail templates: %v", err)
		}
	}

	var mailer mail.Mailer
	if cfg.MailBackend == "console" {
		mailer = mail.NewConsoleMailer(tmpl)
	} else {
		mailer = mail.NewSMTPMailer(mail.SMTPConfig{
			Host:     cfg.SMTPHost,
			Port:     cfg.SMTPPort,
			Username: cfg.SMTPUsername,
			Password: cfg.SMTPPassword,
			From:     cfg.SMTPFrom,
			FromName: cfg.SMTPFromName,
			TLS:      cfg.SMTPTLS,
		}, tmpl)
	}

	users := postgres.NewUserRepo(pool)
	methods := postgres.NewAuthMethodRepo(pool)
	sessions := postgres.NewSessionRepo(pool)
	settings := postgres.NewSettingsRepo(pool)
	codes := postgres.NewMagicCodeRepo(pool)
	audit := postgres.NewLoginAuditRepo(pool)
	externalAccounts := postgres.NewExternalSiteAccountRepo(pool)

	credentialCipher, err := credentials.NewCipher(cfg.ExternalPasswordEncryptionKey)
	if err != nil {
		log.Fatalf("external credentials: %v", err)
	}
	accountFactory, err := credentials.NewAccountFactory(credentialCipher, []entity.ExternalSite{entity.ExternalSiteTPAI, entity.ExternalSiteAPAI})
	if err != nil {
		log.Fatalf("external accounts: %v", err)
	}
	tpaiClient, err := newapi.NewClient(cfg.TPAI.BaseURL, cfg.TPAI.AdminToken, 10*time.Second)
	if err != nil {
		log.Fatalf("TPAI client: %v", err)
	}
	apaiClient, err := newapi.NewClient(cfg.APAI.BaseURL, cfg.APAI.AdminToken, 10*time.Second)
	if err != nil {
		log.Fatalf("APAI client: %v", err)
	}

	authUC := usecase.NewAuthUsecase(
		users, methods, sessions, settings, codes, audit, mailer, geoip.NewHTTPResolver(),
		usecase.WithNewAccountProvisioning(externalAccounts, accountFactory),
	)
	profileUC := usecase.NewProfileUsecase(users, methods, sessions, settings)
	settingsUC := usecase.NewSettingsUsecase(settings, users)

	google := oauth.NewGoogle(oauth.ProviderConfig{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
	})
	github := oauth.NewGithub(oauth.ProviderConfig{
		ClientID:     cfg.GithubClientID,
		ClientSecret: cfg.GithubClientSecret,
		RedirectURL:  cfg.GithubRedirectURL,
	})

	sessionsMW := middleware.NewSessionAuth(authUC, cfg.SessionCookieName, cfg.SessionTTLDays)
	authH := controller.NewAuthController(authUC, sessionsMW, google, github, cfg.PublicWebURL)
	profileH := controller.NewProfileController(profileUC, authUC, google, github)
	settingsH := controller.NewSettingsController(settingsUC)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true, "time": time.Now().UTC()})
	})

	v1 := r.Group("/api/v1")
	v1.Use(middleware.RequireCSRF())
	{
		v1.POST("/auth/magic/request", authH.MagicRequest)
		v1.POST("/auth/magic/resend", authH.MagicResend)
		v1.POST("/auth/magic/verify", authH.MagicVerify)
		v1.GET("/auth/oauth/:provider", authH.OAuthStart)
		v1.GET("/auth/oauth/:provider/callback", authH.OAuthCallback)

		authed := v1.Group("")
		authed.Use(sessionsMW.RequireSession())
		{
			authed.GET("/auth/me", authH.Me)
			authed.POST("/auth/logout", authH.Logout)

			authed.GET("/profile", profileH.Get)
			authed.PATCH("/profile", profileH.Patch)
			authed.GET("/profile/sessions", profileH.ListSessions)
			authed.DELETE("/profile/sessions/:id", profileH.RevokeSession)
			authed.POST("/profile/sessions/revoke-all", profileH.RevokeAllSessions)
			authed.GET("/profile/connections", profileH.ListConnections)
			authed.POST("/profile/connections/:provider/link", profileH.LinkStart)
			authed.DELETE("/profile/connections/:provider", profileH.Unlink)
			authed.GET("/profile/security", profileH.GetSecurity)
			authed.PATCH("/profile/security", profileH.PatchSecurity)

			authed.GET("/settings", settingsH.Get)
			authed.PATCH("/settings", settingsH.Patch)
		}
	}

	externalWorker, err := worker.NewExternalSiteAccountWorker(externalAccounts, credentialCipher, map[entity.ExternalSite]worker.UserCreator{
		entity.ExternalSiteTPAI: tpaiClient,
		entity.ExternalSiteAPAI: apaiClient,
	})
	if err != nil {
		log.Fatalf("external account worker: %v", err)
	}

	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	workerDone := externalWorker.Start(appCtx)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server: %v", err)
		}
	case <-appCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown: %v", err)
		}
		cancel()
	}
	stop()
	<-workerDone
}

func runMigrations(databaseURL string) error {
	if os.Getenv("SKIP_MIGRATE") == "1" {
		return nil
	}
	src, err := (&file.File{}).Open("file://migrations")
	if err != nil {
		return fmt.Errorf("open migrations: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("file", src, databaseURL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
