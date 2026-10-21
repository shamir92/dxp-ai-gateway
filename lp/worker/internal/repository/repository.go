package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/worker/internal/entity"
)

var ErrNotFound = errors.New("repository: not found")
var ErrLeaseLost = errors.New("repository: external account lease lost")

type ExternalSiteAccountRepo interface {
	ClaimDue(ctx context.Context, now, leaseExpiresAt time.Time, limit int) ([]entity.ExternalSiteAccount, error)
	MarkProvisioned(ctx context.Context, id, leaseToken uuid.UUID, remoteUserID string, at time.Time) error
	MarkRetry(ctx context.Context, id, leaseToken uuid.UUID, nextAttemptAt time.Time, summary string, at time.Time) error
	MarkFailed(ctx context.Context, id, leaseToken uuid.UUID, summary string, at time.Time) error
}
