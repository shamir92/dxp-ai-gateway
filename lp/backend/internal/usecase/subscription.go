package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/newapi"
)

// QuotaToCreditsFactor converts upstream quota units into the credits shown on
// our side. Both amount_total and amount_used are scaled by it.
const QuotaToCreditsFactor int64 = 1428

// SubscriptionView is the DXP-facing shape of one site's subscription usage. It
// deliberately does not name the backing platform.
type SubscriptionView struct {
	PlanID            int    `json:"plan_id"`
	Status            string `json:"status"`
	AmountTotal       int64  `json:"amount_total"`
	AmountUsed        int64  `json:"amount_used"`
	CreditsTotal      *int64 `json:"credits_total"`
	CreditsUsed       int64  `json:"credits_used"`
	Unlimited         bool   `json:"unlimited"`
	StartTime         int64  `json:"start_time"`
	EndTime           int64  `json:"end_time"`
	NextResetTime     int64  `json:"next_reset_time"`
	BillingPreference string `json:"billing_preference"`
}

// Endpoint is one API address advertised by the platform. The same base is used
// for OpenAI-compatible and Anthropic-compatible clients.
type Endpoint struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Endpoints lists the platform's advertised API addresses, taken from the
// public GET /api/status (data.api_info). A platform may advertise none.
func (u *APIKeyUsecase) Endpoints(ctx context.Context, userID uuid.UUID, site entity.ExternalSite) ([]Endpoint, error) {
	client, err := u.clientFor(site)
	if err != nil {
		return nil, err
	}
	status, err := client.Status(ctx)
	if err != nil {
		return nil, err
	}
	endpoints := make([]Endpoint, 0, len(status.APIInfo))
	for _, info := range status.APIInfo {
		url := strings.TrimRight(strings.TrimSpace(info.URL), "/")
		if url == "" {
			continue
		}
		label := strings.TrimSpace(info.Description)
		if label == "" {
			label = strings.TrimSpace(info.Route)
		}
		endpoints = append(endpoints, Endpoint{Label: label, URL: url})
	}
	return endpoints, nil
}

// Subscription reads the user's subscription usage for one site.
func (u *APIKeyUsecase) Subscription(ctx context.Context, userID uuid.UUID, site entity.ExternalSite) (*SubscriptionView, error) {
	return withBearer(u, ctx, userID, site, func(bearer string) (*SubscriptionView, error) {
		client, err := u.clientFor(site)
		if err != nil {
			return nil, err
		}
		self, err := client.SubscriptionSelf(ctx, bearer)
		if err != nil {
			return nil, err
		}
		return buildSubscriptionView(self), nil
	})
}

// buildSubscriptionView maps the upstream payload. The active plan is the first
// entry in Subscriptions; a total of zero means unlimited, so credits_total stays
// nil in that case rather than being reported as 0.
func buildSubscriptionView(self newapi.SubscriptionSelf) *SubscriptionView {
	view := &SubscriptionView{BillingPreference: self.BillingPreference}
	if len(self.Subscriptions) == 0 {
		return view
	}

	sub := self.Subscriptions[0].Subscription
	view.PlanID = sub.PlanID
	view.Status = sub.Status
	view.AmountTotal = sub.AmountTotal
	view.AmountUsed = sub.AmountUsed
	view.StartTime = sub.StartTime
	view.EndTime = sub.EndTime
	view.NextResetTime = sub.NextResetTime
	view.Unlimited = sub.AmountTotal == 0
	view.CreditsUsed = sub.AmountUsed * QuotaToCreditsFactor
	if !view.Unlimited {
		total := sub.AmountTotal * QuotaToCreditsFactor
		view.CreditsTotal = &total
	}
	return view
}
