package main

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/worker/internal/entity"
	"github.com/dxpglobal/ai-gateway/worker/internal/newapi"
	"github.com/dxpglobal/ai-gateway/worker/internal/repository"
)

const (
	defaultPollInterval = 5 * time.Second
	defaultLease        = time.Minute
	defaultBatchSize    = 2
	maximumRetryDelay   = 6 * time.Hour
)

type PasswordDecryptor interface {
	Decrypt(userID uuid.UUID, site entity.ExternalSite, ciphertext, nonce []byte, keyVersion string) ([]byte, error)
}

type UserCreator interface {
	CreateUser(ctx context.Context, username, displayName string, password []byte) error
}

// Provisioner polls the external_site_accounts queue and creates the matching
// remote account on each external platform.
type Provisioner struct {
	repo         repository.ExternalSiteAccountRepo
	decryptor    PasswordDecryptor
	clients      map[entity.ExternalSite]UserCreator
	pollInterval time.Duration
	lease        time.Duration
	batchSize    int
}

func NewProvisioner(repo repository.ExternalSiteAccountRepo, decryptor PasswordDecryptor, clients map[entity.ExternalSite]UserCreator) (*Provisioner, error) {
	if repo == nil || decryptor == nil || clients[entity.ExternalSiteTPAI] == nil || clients[entity.ExternalSiteAPAI] == nil {
		return nil, errors.New("provisioner dependencies are incomplete")
	}
	return &Provisioner{
		repo:         repo,
		decryptor:    decryptor,
		clients:      clients,
		pollInterval: defaultPollInterval,
		lease:        defaultLease,
		batchSize:    defaultBatchSize,
	}, nil
}

func (p *Provisioner) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(p.pollInterval)
		defer ticker.Stop()
		for {
			p.processDue(ctx)
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
	return done
}

func (p *Provisioner) processDue(ctx context.Context) {
	now := time.Now().UTC()
	accounts, err := p.repo.ClaimDue(ctx, now, now.Add(p.lease), p.batchSize)
	if err != nil {
		if ctx.Err() == nil {
			// The error is the only thing that explains a stalled queue; the
			// previous message discarded it.
			slog.Error("claim failed", slog.Any("error", err))
			sentry.CaptureException(err)
		}
		return
	}
	for _, account := range accounts {
		if ctx.Err() != nil {
			return
		}
		p.processAccount(ctx, account)
	}
}

func (p *Provisioner) processAccount(ctx context.Context, account entity.ExternalSiteAccount) {
	client := p.clients[account.Site]
	if client == nil {
		p.fail(ctx, account, "unsupported external site")
		return
	}
	password, err := p.decryptor.Decrypt(account.UserID, account.Site, account.PasswordCiphertext, account.PasswordNonce, account.EncryptionKeyVersion)
	if err != nil {
		p.fail(ctx, account, "stored credential could not be decrypted")
		return
	}
	err = client.CreateUser(ctx, account.Username, account.DisplayName, password)
	clear(password)
	if err == nil {
		if err := p.repo.MarkProvisioned(ctx, account.ID, account.LeaseToken, "", time.Now().UTC()); err != nil {
			slog.Error("completion update failed",
				slog.String("account_id", account.ID.String()),
				slog.String("site", string(account.Site)),
				slog.Any("error", err),
			)
			sentry.CaptureException(err)
		}
		return
	}

	var apiErr *newapi.APIError
	if errors.As(err, &apiErr) && !apiErr.Retryable {
		p.fail(ctx, account, apiErr.Error())
		return
	}
	delay := retryDelay(account.AttemptCount + 1)
	if errors.As(err, &apiErr) && apiErr.RetryAfter > 0 {
		delay = apiErr.RetryAfter
	}
	now := time.Now().UTC()
	slog.Warn("provisioning attempt failed, retrying",
		slog.String("account_id", account.ID.String()),
		slog.String("site", string(account.Site)),
		slog.Int("attempt", account.AttemptCount+1),
		slog.Duration("retry_in", delay),
		slog.String("reason", safeErrorSummary(err)),
	)
	if err := p.repo.MarkRetry(ctx, account.ID, account.LeaseToken, now.Add(delay), safeErrorSummary(err), now); err != nil {
		slog.Error("retry update failed",
			slog.String("account_id", account.ID.String()),
			slog.Any("error", err),
		)
		sentry.CaptureException(err)
	}
}

func (p *Provisioner) fail(ctx context.Context, account entity.ExternalSiteAccount, summary string) {
	// A failed account is the worker's real incident signal: nothing else
	// surfaces it, and it means a user is stuck mid-provisioning.
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("site", string(account.Site))
		scope.SetContext("account", map[string]any{
			"account_id": account.ID.String(),
			"attempts":   account.AttemptCount,
		})
		sentry.CaptureMessage("external account provisioning failed: " + summary)
	})
	slog.Error("provisioning failed permanently",
		slog.String("account_id", account.ID.String()),
		slog.String("site", string(account.Site)),
		slog.Int("attempts", account.AttemptCount+1),
		slog.String("reason", summary),
	)
	if err := p.repo.MarkFailed(ctx, account.ID, account.LeaseToken, summary, time.Now().UTC()); err != nil {
		slog.Error("failure update failed",
			slog.String("account_id", account.ID.String()),
			slog.Any("error", err),
		)
		sentry.CaptureException(err)
	}
}

func safeErrorSummary(err error) string {
	var apiErr *newapi.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Error()
	}
	return "external service request failed"
}

func retryDelay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	shift := min(attempt-1, 9)
	base := time.Minute * time.Duration(1<<shift)
	if base > maximumRetryDelay {
		base = maximumRetryDelay
	}
	jitter := 0.75 + rand.Float64()*0.5
	return time.Duration(float64(base) * jitter)
}
