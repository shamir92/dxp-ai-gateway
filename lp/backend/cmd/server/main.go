package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/backend/internal/controller"
	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/credentials"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/newapi"
	"github.com/dxpglobal/ai-gateway/backend/internal/middleware"
	"github.com/dxpglobal/ai-gateway/backend/internal/migrations"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/config"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/geoip"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/logging"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/mail"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/oauth"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/observability"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository/postgres"
	"github.com/dxpglobal/ai-gateway/backend/internal/usecase"
)

func main() {
	// Load lp/backend/.env before config (process env wins over file).
	if err := config.LoadDotEnv(".env"); err != nil {
		if !os.IsNotExist(err) {
			slog.Warn("dotenv load failed", slog.Any("error", err))
		} else if err := config.LoadDotEnv("lp/backend/.env"); err != nil && !os.IsNotExist(err) {
			slog.Warn("dotenv load failed", slog.Any("error", err))
		}
	}

	cfg, err := config.Load()
	if err != nil {
		fatal("config", err)
	}

	// Install structured logging once configuration (and therefore LOG_LEVEL
	// from .env) is known.
	logging.Setup()

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		fatal("migrate", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		fatal("db", err)
	}
	defer pool.Close()

	tmpl, err := mail.LoadTemplates()
	if err != nil {
		fatal("mail templates", err)
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
		fatal("external credentials", err)
	}
	accountFactory, err := credentials.NewAccountFactory(credentialCipher, []entity.ExternalSite{entity.ExternalSiteTPAI, entity.ExternalSiteAPAI})
	if err != nil {
		fatal("external accounts", err)
	}
	tpaiClient, err := newapi.NewClient(cfg.TPAI.BaseURL, cfg.TPAI.AdminToken, cfg.CloudflareAccess.ClientID, cfg.CloudflareAccess.ClientSecret, 10*time.Second)
	if err != nil {
		fatal("TPAI client", err)
	}
	apaiClient, err := newapi.NewClient(cfg.APAI.BaseURL, cfg.APAI.AdminToken, cfg.CloudflareAccess.ClientID, cfg.CloudflareAccess.ClientSecret, 10*time.Second)
	if err != nil {
		fatal("APAI client", err)
	}

	authUC := usecase.NewAuthUsecase(
		users, methods, sessions, settings, codes, audit, mailer, geoip.NewHTTPResolver(),
		usecase.WithNewAccountProvisioning(externalAccounts, accountFactory),
	)
	profileUC := usecase.NewProfileUsecase(users, methods, sessions, settings)
	settingsUC := usecase.NewSettingsUsecase(settings, users)
	apiKeyUC := usecase.NewAPIKeyUsecase(users, externalAccounts, externalAccounts, credentialCipher, accountFactory, map[entity.ExternalSite]usecase.SiteClient{
		entity.ExternalSiteTPAI: tpaiClient,
		entity.ExternalSiteAPAI: apaiClient,
	})

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
	profileH := controller.NewProfileController(profileUC, authUC, google, github, apiKeyUC)
	settingsH := controller.NewSettingsController(settingsUC)
	apiKeysH := controller.NewAPIKeyController(apiKeyUC)

	// Sentry is optional: without SENTRY_DSN this is a no-op.
	observability.Init()
	defer observability.Flush()

	// gin.New rather than gin.Default: the middleware chain is explicit, and
	// request logging is ours (structured, correlated) instead of gin's.
	r := gin.New()
	r.Use(
		gin.Recovery(),
		middleware.RequestID(),
		sentrygin.New(sentrygin.Options{
			// Repanic so gin.Recovery writes the response; sentrygin captures
			// the event on the way past.
			Repanic: true,
		}),
		middleware.RequestLogger(),
	)
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

			// API keys and subscription usage are proxied with one stored per-user
			// access token. The path uses a neutral product slug so the upstream
			// platform is never named in a URL: tp_ = token plan.
			authed.GET("/profile/tp_api_keys", controller.WithSite(entity.ExternalSiteTPAI), apiKeysH.List)
			authed.POST("/profile/tp_api_keys", controller.WithSite(entity.ExternalSiteTPAI), apiKeysH.Create)
			authed.PUT("/profile/tp_api_keys/:id", controller.WithSite(entity.ExternalSiteTPAI), apiKeysH.Update)
			authed.DELETE("/profile/tp_api_keys/:id", controller.WithSite(entity.ExternalSiteTPAI), apiKeysH.Delete)
			authed.GET("/profile/tp_subscription", controller.WithSite(entity.ExternalSiteTPAI), apiKeysH.Subscription)
			authed.GET("/profile/tp_endpoints", controller.WithSite(entity.ExternalSiteTPAI), apiKeysH.Endpoints)
		}
	}

	// External account provisioning runs in cmd/worker, deployed separately.
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("listening", slog.String("addr", cfg.HTTPAddr))
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server stopped", slog.Any("error", err))
		}
	case <-appCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("server shutdown failed", slog.Any("error", err))
		}
		cancel()
	}
	stop()
}

func runMigrations(databaseURL string) error {
	if os.Getenv("SKIP_MIGRATE") == "1" {
		return nil
	}
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("open migrations: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, databaseURL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// fatal reports a startup failure as structured JSON and exits non-zero, so a
// container platform sees the failure instead of a running-but-unusable server.
func fatal(what string, err error) {
	slog.Error("startup failed", slog.String("stage", what), slog.Any("error", err))
	os.Exit(1)
}
