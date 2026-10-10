package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type Profile struct {
	ProviderUserID string
	Email          string
	Name           string
	AvatarURL      string
}

type ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type GoogleProvider struct {
	cfg  ProviderConfig
	oauth *oauth2.Config
}

func NewGoogle(cfg ProviderConfig) *GoogleProvider {
	return &GoogleProvider{
		cfg: cfg,
		oauth: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (p *GoogleProvider) AuthCodeURL(state string) string {
	return p.oauth.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (p *GoogleProvider) Exchange(ctx context.Context, code string) (*Profile, error) {
	tok, err := p.oauth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("google exchange: %w", err)
	}
	client := p.oauth.Client(ctx, tok)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return nil, fmt.Errorf("google userinfo: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Sub     string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if raw.Email == "" {
		return nil, fmt.Errorf("google: no email")
	}
	return &Profile{
		ProviderUserID: raw.Sub,
		Email:          raw.Email,
		Name:           raw.Name,
		AvatarURL:      raw.Picture,
	}, nil
}

type GithubProvider struct {
	cfg   ProviderConfig
	oauth *oauth2.Config
}

func NewGithub(cfg ProviderConfig) *GithubProvider {
	return &GithubProvider{
		cfg: cfg,
		oauth: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (p *GithubProvider) AuthCodeURL(state string) string {
	return p.oauth.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (p *GithubProvider) Exchange(ctx context.Context, code string) (*Profile, error) {
	tok, err := p.oauth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("github exchange: %w", err)
	}
	client := p.oauth.Client(ctx, tok)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("github user: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var raw struct {
		ID        json.Number `json:"id"`
		Login     string      `json:"login"`
		Name      string      `json:"name"`
		Email     string      `json:"email"`
		AvatarURL string      `json:"avatar_url"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	email := raw.Email
	if email == "" {
		email, err = fetchGithubPrimaryEmail(ctx, client)
		if err != nil {
			return nil, err
		}
	}
	name := raw.Name
	if name == "" {
		name = raw.Login
	}
	return &Profile{
		ProviderUserID: raw.ID.String(),
		Email:          email,
		Name:           name,
		AvatarURL:      raw.AvatarURL,
	}, nil
}

func fetchGithubPrimaryEmail(ctx context.Context, client *http.Client) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary && e.Email != "" {
			return e.Email, nil
		}
	}
	if len(emails) > 0 && emails[0].Email != "" {
		return emails[0].Email, nil
	}
	return "", fmt.Errorf("github: no email")
}
