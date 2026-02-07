# Template

Шаблон Go-сервиса: gRPC + HTTP gateway, PostgreSQL, OpenTelemetry (трейсы + метрики), Temporal workflows.

## Требования

- Go 1.24+
- Docker
- [Task](https://taskfile.dev) (`go install github.com/go-task/task/v3/cmd/task@latest`)

## Быстрый старт

```bash
# 1. Установить инструменты (buf, goose, configgen, protogen и т.д.)
task install

# 2. Создать .env в deployments/local/
cp deployments/local/.env.example deployments/local/.env

# 3. Поднять PostgreSQL + Temporal
task deps

# 4. Применить миграции
task migrate

# 5. Запустить приложение
task run
```

После запуска доступны:

| Сервис | Адрес |
|--------|-------|
| HTTP gateway | `http://localhost:7001` |
| gRPC | `localhost:7000` |
| Swagger UI | `http://localhost:7002` |
| Debug (pprof) | `http://localhost:7003/debug/pprof/` |
| Health check | `http://localhost:7003/healthz` |
| Prometheus metrics | `http://localhost:7003/metrics` |

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
├── api/                    # Proto-файлы и buf конфигурация
├── cmd/template/           # Точка входа: main.go, dependency.go
├── config/                 # TOML конфиги + сгенерированные Go-структуры
├── deployments/local/      # Docker Compose, Prometheus, Grafana, Tempo
├── internal/
│   ├── controller/         # gRPC контроллеры
│   └── pkg/
│       ├── metrics/        # Prometheus handler
│       ├── otel/           # Инициализация OpenTelemetry
│       ├── storage/postgres/ # PgX пул + транзакции
│       └── ...
├── migrations/             # SQL миграции (goose)
└── pkg/                    # Сгенерированный proto-код + Swagger embed
```

## Конфигурация

Используется `configgen` — генерирует Go-структуры из TOML-файлов.

- `config/value.toml` — общие значения для всех окружений
- `config/config_local.toml` — настройки для локальной разработки
- `config/config_prod.toml` — настройки для production

Окружение выбирается через `APP_ENV` (default: `dev`).

```bash
# Перегенерировать Go-структуры после изменения TOML
task generate-config
```

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
| `task proto:gen` | Генерация Go-кода из proto |
| `task proto:controllers` | Генерация stub-контроллеров |
| `task generate-config` | Генерация Go-структур конфига |
| `task generate` | Все генераторы разом |

### База данных
| Команда | Описание |
|---------|----------|
| `task migrate` | Применить миграции |
| `task migrate:down` | Откатить одну миграцию |
| `task migrate:create -- name` | Создать новую миграцию |

### Инфраструктура
| Команда | Описание |
|---------|----------|
| `task deps` | Запустить PostgreSQL + Temporal |
| `task deps:stop` | Остановить |
| `task deps:logs` | Логи |
| `task metrics` | Запустить Prometheus + Grafana + Tempo + OTEL Collector |
| `task metrics:stop` | Остановить стек наблюдаемости |

## Технологии

- **Сервер**: [platform](https://github.com/vovanwin/platform) (gRPC + HTTP gateway + Swagger + Debug)
- **DI**: uber/fx
- **БД**: PostgreSQL (pgx/v5), миграции goose
- **Proto**: buf, grpc-gateway
- **Observability**: OpenTelemetry, Prometheus, Grafana, Tempo
- **Workflows**: Temporal
- **Auth**: JWT (golang-jwt), Argon2ID
