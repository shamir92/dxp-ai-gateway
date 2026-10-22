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

// SubscriptionView is the DXP-facing shape of one plan instance. It deliberately
// does not name the backing platform.
type SubscriptionView struct {
	// SubscriptionID identifies the plan instance. It is not unique per plan_id:
	// a user can hold several instances of the same plan.
	SubscriptionID int    `json:"subscription_id"`
	PlanID         int    `json:"plan_id"`
	Status         string `json:"status"`
	AmountTotal    int64  `json:"amount_total"`
	AmountUsed     int64  `json:"amount_used"`
	CreditsTotal   *int64 `json:"credits_total"`
	CreditsUsed    int64  `json:"credits_used"`
	Unlimited      bool   `json:"unlimited"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	NextResetTime  int64  `json:"next_reset_time"`
}

// SubscriptionUsage is the whole of a user's subscription state for one site.
// The platform lets a user hold more than one active plan, so plans are a list
// rather than a single record.
type SubscriptionUsage struct {
	BillingPreference string             `json:"billing_preference"`
	Plans             []SubscriptionView `json:"plans"`
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
func (u *APIKeyUsecase) Subscription(ctx context.Context, userID uuid.UUID, site entity.ExternalSite) (*SubscriptionUsage, error) {
	return withBearer(u, ctx, userID, site, func(bearer string) (*SubscriptionUsage, error) {
		client, err := u.clientFor(site)
		if err != nil {
			return nil, err
		}
		self, err := client.SubscriptionSelf(ctx, bearer)
		if err != nil {
			return nil, err
		}
		return buildSubscriptionUsage(self), nil
	})
}

// buildSubscriptionUsage maps the upstream payload. Plans keeps every active
// subscription in the order the platform returned them, which is the plan
// priority order. A total of zero means unlimited, so credits_total stays nil
// in that case rather than being reported as 0.
func buildSubscriptionUsage(self newapi.SubscriptionSelf) *SubscriptionUsage {
	usage := &SubscriptionUsage{
		BillingPreference: self.BillingPreference,
		Plans:             make([]SubscriptionView, 0, len(self.Subscriptions)),
	}
	for _, summary := range self.Subscriptions {
		usage.Plans = append(usage.Plans, buildSubscriptionView(summary.Subscription))
	}
	return usage
}

func buildSubscriptionView(sub newapi.Subscription) SubscriptionView {
	view := SubscriptionView{
		SubscriptionID: sub.ID,
		PlanID:         sub.PlanID,
		Status:         sub.Status,
		AmountTotal:    sub.AmountTotal,
		AmountUsed:     sub.AmountUsed,
		StartTime:      sub.StartTime,
		EndTime:        sub.EndTime,
		NextResetTime:  sub.NextResetTime,
		Unlimited:      sub.AmountTotal == 0,
		CreditsUsed:    sub.AmountUsed * QuotaToCreditsFactor,
	}
	if !view.Unlimited {
		total := sub.AmountTotal * QuotaToCreditsFactor
		view.CreditsTotal = &total
	}
	return view
}
