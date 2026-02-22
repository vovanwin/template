# Template

Шаблон Go-сервиса: gRPC + HTTP gateway + Web UI (HTMX + Templ), PostgreSQL, Centrifugo (realtime), Telegram-бот, OpenTelemetry (трейсы + метрики), Temporal workflows.

## Возможности

- **gRPC + HTTP Gateway** — единый API через proto-определения, автогенерация REST-эндпоинтов через grpc-gateway
- **Web UI** — серверный рендеринг на Templ + интерактивность через HTMX, Alpine.js для состояния на клиенте
- **Realtime-уведомления** — Centrifugo (uni_sse) с JWT-авторизацией, персональные каналы, демо-страница
- **Telegram-бот** — интеграция с Telegram Bot API, поддержка Mini App (WebApp)
- **Напоминания** — CRUD + выполнение через Temporal Workflows с подтверждением и повторами
- **Аутентификация** — JWT (access + refresh), Argon2ID хеширование паролей, CSRF-защита
- **Feature Flags** — etcd-хранилище с UI на debug-порту, горячая перезагрузка
- **Observability** — OpenTelemetry (трейсы + метрики), Prometheus, Grafana, Tempo, Loki
- **Конфигурация** — `configgen` генерирует типизированные Go-структуры из TOML, мультиокружения
- **Компонентный логгер** — slog с переопределением уровней по компонентам, отправка в Loki

## Требования

- Go 1.24+
- Docker
- [Task](https://taskfile.dev) (`go install github.com/go-task/task/v3/cmd/task@latest`)

## Быстрый старт

```bash
# 1. Установить инструменты (buf, goose, templ, configgen, protogen и т.д.)
task install

# 2. Создать .env файл для docker-compose
cp .env_example .env

# 3. Поднять PostgreSQL + Temporal + Centrifugo
task deps

# 4. Применить миграции
task migrate

# 5. Запустить приложение
task run
```

После запуска доступны:

| Сервис | Адрес |
|--------|-------|
| Web UI | `http://localhost:7001` |
| gRPC | `localhost:7000` |
| Swagger UI | `http://localhost:7002` |
| Debug (pprof) | `http://localhost:7003/debug/pprof/` |
| Health check | `http://localhost:7003/healthz` |
| Prometheus metrics | `http://localhost:7003/metrics` |
| Feature Flags UI | `http://localhost:7003/flags` |
| Centrifugo Admin | `http://localhost:8000` (password/admin_secret) |

## Realtime-уведомления (Centrifugo)

Архитектура:

```
Browser ──EventSource (uni_sse)──→ Centrifugo ←── HTTP API ── Go Backend
                                       │
                                  JWT авторизация
                                  personal#<userID>
```

- Бэкенд генерирует JWT токен → клиент подключается к Centrifugo через нативный `EventSource`
- Автоподписка на персональный канал `#<userID>`
- Бэкенд публикует события через HTTP Server API (`POST /api/publish`)
- Типы уведомлений: success, info, warning, error, reminder — с цветовой индикацией
- Демо-страница: `/notifications-demo` — интерактивная демонстрация всех типов

```bash
# Запуск Centrifugo отдельно (если не в основном docker-compose)
docker compose -f deployments/local/docker-compose.centrifugo.yml up -d
```

Конфигурация в `config/config_local.toml`:

```toml
[centrifugo]
addr = "http://localhost:8000"
api_key = "centrifugo-api-key"
token_secret = "centrifugo-token-secret"
token_ttl = "60m"
```

## Observability (Prometheus + Grafana + Tempo)

```bash
# Поднять стек наблюдаемости
task metrics

# Остановить
task metrics:stop

# Логи
task metrics:logs
```

| Сервис | Адрес |
|--------|-------|
| Grafana | `http://localhost:3000` (admin/admin) |
| Prometheus | `http://localhost:9090` |
| Tempo | `http://localhost:3200` |
| OTEL Collector | `localhost:4317` (gRPC) |

Трейсы и метрики отправляются через OTLP gRPC в OTEL Collector, который маршрутизирует их в Tempo (трейсы) и Prometheus (метрики). Приложение также экспонирует `/metrics` на debug-порту для прямого скрейпинга Prometheus.

Управляется через конфиг:

```toml
# config/value.toml
[features]
enable_metrics = true
enable_tracing = true

[otel]
endpoint = "localhost:4317"
```

## Структура проекта

```
├── api/                    # Proto-файлы для gRPC API (buf)
├── api-workflow/           # Proto-файлы для Temporal Workflows
├── cmd/template/           # Точка входа: main.go, dependency.go
├── config/                 # TOML конфиги + сгенерированные Go-структуры
├── deployments/local/      # Docker Compose + Centrifugo config
├── internal/
│   ├── controller/         # gRPC и HTTP контроллеры
│   │   ├── auth/           # Контроллеры аутентификации
│   │   ├── ui/             # Web UI (HTMX/Templ): pages, layouts, components
│   │   └── ...
│   ├── pkg/
│   │   ├── centrifugo/     # HTTP-клиент Centrifugo Server API + JWT
│   │   ├── events/         # Event bus (публикация через Centrifugo)
│   │   ├── jwt/            # JWT-сервис (access/refresh токены)
│   │   ├── telegram/       # Telegram-бот (модульная архитектура)
│   │   └── ...
│   ├── repository/         # Слой доступа к данным (PostgreSQL)
│   ├── service/            # Бизнес-логика
│   └── workflows/          # Temporal Workflows и Activities
├── migrations/             # SQL миграции (goose)
└── pkg/                    # Сгенерированный proto-код + Swagger embed
```

## Конфигурация

Используется `configgen` — генерирует Go-структуры из TOML-файлов.

- `config/value.toml` — общие значения для всех окружений
- `config/config_local.toml` — настройки для локальной разработки
- `config/config_prod.toml` — настройки для production
- `config/flags.toml` — feature flags с дефолтными значениями

Окружение выбирается через `APP_ENV` (default: `local`).

```bash
# Перегенерировать Go-структуры и флаги после изменения TOML
task generate-config

# Проверить валидность конфигов без генерации
task validate-config
```

### Feature Flags

Feature flags определяются в `config/flags.toml` и доступны через типизированные геттеры:

```go
flags := config.NewFlags(config.NewMemoryStore(config.DefaultFlagValues()))

if flags.NewCatalogUi() {
    // новый UI
}
limit := flags.RateLimit() // int
```

UI для просмотра флагов доступен на debug-порту: `http://localhost:7003/flags`

## Task команды

### Приложение
| Команда | Описание |
|---------|----------|
| `task run` | Сборка и запуск |
| `task format` | Форматирование (gofumpt + gci) |
| `task lint` | Линтинг (golangci-lint) |

### Кодогенерация
| Команда | Описание |
|---------|----------|
| `task generate` | Все генераторы разом (proto, templ, config) |
| `task proto:gen` | Генерация Go-кода из proto (gRPC и Temporal) |
| `task templ:gen` | Генерация Go-кода из .templ файлов |
| `task generate-config` | Генерация Go-структур конфига + feature flags |
| `task proto:controllers` | Генерация stub-контроллеров |
| `task validate-config` | Валидация TOML без генерации (для CI) |


### База данных
| Команда | Описание |
|---------|----------|
| `task migrate` | Применить миграции |
| `task migrate:down` | Откатить одну миграцию |
| `task migrate:create -- name` | Создать новую миграцию |

### Инфраструктура
| Команда | Описание |
|---------|----------|
| `task deps` | Запустить PostgreSQL + Temporal + Centrifugo |
| `task deps:stop` | Остановить |
| `task deps:logs` | Логи |
| `task metrics` | Запустить Prometheus + Grafana + Tempo + OTEL Collector |
| `task metrics:stop` | Остановить стек наблюдаемости |

## Технологии

| Категория | Стек |
|-----------|------|
| Сервер | [platform](https://github.com/vovanwin/platform) (gRPC + HTTP gateway + Swagger + Debug) |
| Web UI | [HTMX](https://htmx.org/) + [Templ](https://templ.guide/) + [Alpine.js](https://alpinejs.dev/) + Tailwind CSS |
| Realtime | [Centrifugo](https://centrifugal.dev/) v5 (uni_sse, JWT, персональные каналы) |
| Telegram | [go-telegram/bot](https://github.com/go-telegram/bot) + Mini App (WebApp) |
| DI | [uber/fx](https://github.com/uber-go/fx) |
| БД | PostgreSQL ([pgx/v5](https://github.com/jackc/pgx)), миграции [goose](https://github.com/pressly/goose) |
| Proto | [buf](https://buf.build/), grpc-gateway, protoc-gen-go-temporal |
| Observability | OpenTelemetry, Prometheus, Grafana, Tempo, Loki |
| Workflows | [Temporal](https://temporal.io/) |
| Auth | JWT ([golang-jwt](https://github.com/golang-jwt/jwt)), Argon2ID, CSRF |
| Config | [configgen](https://github.com/vovanwin/configgen) (TOML → Go structs) |
| Feature Flags | etcd + in-memory fallback, UI на debug-порту |
