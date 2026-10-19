package newapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Subscription is one plan instance for a user. Amounts are in upstream quota
// units; the gateway scales them for display.
type Subscription struct {
	ID                  int    `json:"id"`
	UserID              int    `json:"user_id"`
	PlanID              int    `json:"plan_id"`
	AmountTotal         int64  `json:"amount_total"`
	AmountUsed          int64  `json:"amount_used"`
	StartTime           int64  `json:"start_time"`
	EndTime             int64  `json:"end_time"`
	Status              string `json:"status"`
	Source              string `json:"source"`
	LastResetTime       int64  `json:"last_reset_time"`
	NextResetTime       int64  `json:"next_reset_time"`
	UpgradeGroup        string `json:"upgrade_group"`
	PrevUserGroup       string `json:"prev_user_group"`
	DowngradeGroup      string `json:"downgrade_group"`
	AllowWalletOverflow bool   `json:"allow_wallet_overflow"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
}

// SubscriptionSummary wraps each array entry. The upstream nests the payload at
// entry.subscription rather than putting the fields at the top level.
type SubscriptionSummary struct {
	Subscription Subscription `json:"subscription"`
}

// SubscriptionSelf is the /api/subscription/self payload. Subscriptions holds
// active plans; AllSubscriptions includes expired ones.
type SubscriptionSelf struct {
	BillingPreference string                `json:"billing_preference"`
	Subscriptions     []SubscriptionSummary `json:"subscriptions"`
	AllSubscriptions  []SubscriptionSummary `json:"all_subscriptions"`
}

func (c *Client) SubscriptionSelf(ctx context.Context, bearer string) (SubscriptionSelf, error) {
	data, err := c.do(ctx, requestOptions{
		method:  http.MethodGet,
		path:    subscriptionSelfPath,
		bearer:  bearer,
		referer: "/wallet",
	})
	if err != nil {
		return SubscriptionSelf{}, err
	}

	var result SubscriptionSelf
	if err := json.Unmarshal(data, &result); err != nil {
		return SubscriptionSelf{}, fmt.Errorf("decode subscription response")
	}
	return result, nil
}
