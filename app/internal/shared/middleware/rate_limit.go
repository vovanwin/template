package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// getClientIP извлекает IP адрес клиента из запроса
func getClientIP(r *http.Request) string {
	// Проверяем заголовок X-Forwarded-For
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Берем первый IP из списка
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Проверяем X-Real-IP
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Если заголовки не найдены, используем RemoteAddr
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip, _, _ = net.SplitHostPort(ip)
	}
	return ip
}

// RateLimitRule определяет правило лимитирования
type RateLimitRule struct {
	Requests int           // Количество запросов
	Window   time.Duration // Временное окно
	Path     string        // Путь (если пустой, то для всех)
}

// RateLimiter реализует алгоритм sliding window для ограничения запросов
type RateLimiter struct {
	rules   map[string]RateLimitRule // rules по путям
	clients map[string]*ClientLimit  // состояние клиентов
	mutex   sync.RWMutex
}

// ClientLimit хранит информацию о лимитах клиента
type ClientLimit struct {
	requests []time.Time
	mutex    sync.Mutex
}

// NewRateLimiter создает новый rate limiter
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		rules:   make(map[string]RateLimitRule),
		clients: make(map[string]*ClientLimit),
	}

	// Добавляем правила по умолчанию
	rl.AddRule(
		"/web/login", RateLimitRule{
			Requests: 30,              // 5 попыток входа
			Window:   1 * time.Minute, // за 5 минут
			Path:     "/web/login",
		},
	)

	rl.AddRule(
		"default", RateLimitRule{
			Requests: 300,             // 300 запросов
			Window:   1 * time.Minute, // за минуту
			Path:     "",              // для всех путей
		},
	)

	// Запускаем cleanup горутину
	go rl.cleanup()

	return rl
}

// AddRule добавляет правило лимитирования
func (rl *RateLimiter) AddRule(name string, rule RateLimitRule) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.rules[name] = rule
}

// RateLimitMiddleware возвращает middleware для ограничения запросов
func (rl *RateLimiter) RateLimitMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				clientIP := getClientIP(r)

				// Находим подходящее правило
				rule := rl.findRule(r.URL.Path)

				// Проверяем лимит
				if !rl.isAllowed(clientIP, r.URL.Path, rule) {
					w.Header().Set("Retry-After", fmt.Sprintf("%.0f", rule.Window.Seconds()))
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}

// findRule находит подходящее правило для пути
func (rl *RateLimiter) findRule(path string) RateLimitRule {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	// Ищем специфичное правило для пути
	for _, rule := range rl.rules {
		if rule.Path != "" && rule.Path == path {
			return rule
		}
	}

	// Возвращаем правило по умолчанию
	return rl.rules["default"]
}

// isAllowed проверяет, разрешен ли запрос
func (rl *RateLimiter) isAllowed(clientIP, path string, rule RateLimitRule) bool {
	key := fmt.Sprintf("%s:%s", clientIP, rule.Path)
	if rule.Path == "" {
		key = clientIP // для общих правил используем только IP
	}

	rl.mutex.Lock()
	client, exists := rl.clients[key]
	if !exists {
		client = &ClientLimit{
			requests: make([]time.Time, 0),
		}
		rl.clients[key] = client
	}
	rl.mutex.Unlock()

	client.mutex.Lock()
	defer client.mutex.Unlock()

	now := time.Now()

	// Очищаем старые запросы
	client.requests = filterRecentRequests(client.requests, now, rule.Window)

	// Проверяем лимит
	if len(client.requests) >= rule.Requests {
		return false
	}

	// Добавляем текущий запрос
	client.requests = append(client.requests, now)

	return true
}

// filterRecentRequests оставляет только недавние запросы
func filterRecentRequests(requests []time.Time, now time.Time, window time.Duration) []time.Time {
	cutoff := now.Add(-window)
	result := requests[:0] // повторно используем slice

	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			result = append(result, reqTime)
		}
	}

	return result
}

// cleanup периодически очищает старые записи
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mutex.Lock()
			now := time.Now()

			for key, client := range rl.clients {
				client.mutex.Lock()
				// Если у клиента нет запросов за последний час, удаляем его
				if len(client.requests) == 0 || now.Sub(client.requests[len(client.requests)-1]) > time.Hour {
					delete(rl.clients, key)
				}
				client.mutex.Unlock()
			}

			rl.mutex.Unlock()
		}
	}
}
