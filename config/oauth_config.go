package config

import (
	"os"

	"github.com/vovanwin/template/internal/pkg/oauth"
)

// LoadOAuthProviders загружает конфиги OAuth провайдеров из переменных окружения.
//
// Переменные:
//   - GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, GITHUB_REDIRECT_URL
//   - GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REDIRECT_URL
func LoadOAuthProviders() map[string]oauth.ProviderConfig {
	return map[string]oauth.ProviderConfig{
		"github": {
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		},
		"google": {
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		},
	}
}
