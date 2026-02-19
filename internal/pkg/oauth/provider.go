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

// UserInfo содержит информацию о пользователе от OAuth провайдера.
type UserInfo struct {
	ID    string
	Email string
	Name  string
}

// Provider описывает OAuth провайдера.
type Provider interface {
	Name() string
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*UserInfo, error)
}

// ProviderConfig хранит настройки провайдера.
type ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Registry хранит зарегистрированных провайдеров.
type Registry struct {
	providers map[string]Provider
}

func NewRegistry(configs map[string]ProviderConfig) *Registry {
	r := &Registry{providers: make(map[string]Provider)}

	if cfg, ok := configs["github"]; ok && cfg.ClientID != "" {
		r.providers["github"] = newGitHubProvider(cfg)
	}
	if cfg, ok := configs["google"]; ok && cfg.ClientID != "" {
		r.providers["google"] = newGoogleProvider(cfg)
	}

	return r
}

func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown oauth provider: %s", name)
	}
	return p, nil
}

// ─── GitHub ────────────────────────────────────────────────────────────────

type githubProvider struct {
	cfg *oauth2.Config
}

func newGitHubProvider(c ProviderConfig) Provider {
	return &githubProvider{
		cfg: &oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			RedirectURL:  c.RedirectURL,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"user:email"},
		},
	}
}

func (p *githubProvider) Name() string { return "github" }

func (p *githubProvider) GetAuthURL(state string) string {
	return p.cfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (p *githubProvider) ExchangeCode(ctx context.Context, code string) (*UserInfo, error) {
	token, err := p.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	client := p.cfg.Client(ctx, token)

	// Получаем профиль
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("get github user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github user api returned %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var ghUser struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &ghUser); err != nil {
		return nil, fmt.Errorf("parse github user: %w", err)
	}

	// Email может быть скрыт — если пустой, запрашиваем отдельно
	email := ghUser.Email
	if email == "" {
		email, _ = p.fetchPrimaryEmail(ctx, client)
	}

	return &UserInfo{
		ID:    fmt.Sprintf("%d", ghUser.ID),
		Email: email,
		Name:  ghUser.Name,
	}, nil
}

func (p *githubProvider) fetchPrimaryEmail(ctx context.Context, client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary {
			return e.Email, nil
		}
	}
	if len(emails) > 0 {
		return emails[0].Email, nil
	}
	return "", nil
}

// ─── Google ────────────────────────────────────────────────────────────────

type googleProvider struct {
	cfg *oauth2.Config
}

func newGoogleProvider(c ProviderConfig) Provider {
	return &googleProvider{
		cfg: &oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			RedirectURL:  c.RedirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"openid", "email", "profile"},
		},
	}
}

func (p *googleProvider) Name() string { return "google" }

func (p *googleProvider) GetAuthURL(state string) string {
	return p.cfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (p *googleProvider) ExchangeCode(ctx context.Context, code string) (*UserInfo, error) {
	token, err := p.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	client := p.cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("get google userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo api returned %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var gUser struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &gUser); err != nil {
		return nil, fmt.Errorf("parse google user: %w", err)
	}

	return &UserInfo{
		ID:    gUser.Sub,
		Email: gUser.Email,
		Name:  gUser.Name,
	}, nil
}
