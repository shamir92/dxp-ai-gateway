package geoip

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// Resolver maps an IP to a short place label ("Singapore", "Jakarta, Indonesia").
// Best-effort; never fails the login path.
type Resolver interface {
	Resolve(ctx context.Context, ip string) string
}

type HTTPResolver struct {
	Client *http.Client
}

func NewHTTPResolver() *HTTPResolver {
	return &HTTPResolver{Client: &http.Client{Timeout: 2 * time.Second}}
}

// Resolve uses ip-api.com (free, non-commercial). Falls back to "Unknown location".
func (r *HTTPResolver) Resolve(ctx context.Context, ip string) string {
	ip = strings.TrimSpace(ip)
	if ip == "" || ip == "127.0.0.1" || ip == "::1" || ip == "localhost" {
		return "Local / private network"
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "Unknown location"
	}
	if parsed.IsPrivate() || parsed.IsLoopback() {
		return "Local / private network"
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,regionName,city", ip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "Unknown location"
	}
	resp, err := r.Client.Do(req)
	if err != nil {
		return "Unknown location"
	}
	defer resp.Body.Close()

	var body struct {
		Status    string `json:"status"`
		Country   string `json:"country"`
		Region    string `json:"regionName"`
		City      string `json:"city"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "Unknown location"
	}
	if body.Status != "success" {
		return "Unknown location"
	}
	if body.City != "" && body.Country != "" {
		return body.City + ", " + body.Country
	}
	if body.Country != "" {
		return body.Country
	}
	if body.Region != "" {
		return body.Region
	}
	return "Unknown location"
}

type StaticResolver struct{ Label string }

func (s StaticResolver) Resolve(context.Context, string) string {
	if s.Label == "" {
		return "Unknown location"
	}
	return s.Label
}
