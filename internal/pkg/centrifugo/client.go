package centrifugo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Client struct {
	baseURL     string
	apiKey      string
	tokenSecret string
	tokenTTL    time.Duration
	http        *http.Client
	log         *slog.Logger
}

func NewClient(baseURL, apiKey, tokenSecret string, tokenTTL time.Duration, log *slog.Logger) *Client {
	return &Client{
		baseURL:     baseURL,
		apiKey:      apiKey,
		tokenSecret: tokenSecret,
		tokenTTL:    tokenTTL,
		http:        &http.Client{Timeout: 5 * time.Second},
		log:         log,
	}
}

type publishRequest struct {
	Channel string `json:"channel"`
	Data    any    `json:"data"`
}

// Publish отправляет данные в канал через POST /api/publish.
func (c *Client) Publish(ctx context.Context, channel string, data any) error {
	body, err := json.Marshal(publishRequest{Channel: channel, Data: data})
	if err != nil {
		return fmt.Errorf("centrifugo: marshal publish: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/publish", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("centrifugo: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "apikey "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("centrifugo: publish: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("centrifugo: publish returned status %d", resp.StatusCode)
	}

	c.log.Debug("published to centrifugo", slog.String("channel", channel))
	return nil
}

// PersonalChannel возвращает имя персонального канала для пользователя.
func PersonalChannel(userID string) string {
	return "#" + userID
}

// GenerateToken создаёт JWT для подключения клиента к Centrifugo.
func (c *Client) GenerateToken(userID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(c.tokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.tokenSecret))
}
