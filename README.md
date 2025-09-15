# Шаблон Golang приложения

Минимальный каркас для старта проекта на Go: HTTP (chi), DI (fx), CLI (cobra), миграции (goose), кодоген по OpenAPI (ogen), Postgres (готовая структура для sqlc).

## Обзор проекта

**Template** - это минимальный каркас для старта проекта на Go, предоставляющий готовую архитектуру для создания HTTP REST API с интеграцией различных инструментов и сервисов.

### Ключевые особенности

- **Архитектура**: Модульная архитектура с dependency injection (fx)
- **HTTP Router**: Chi router для обработки HTTP запросов
- **OpenAPI**: Автогенерация кода по OpenAPI спецификации (ogen)
- **База данных**: PostgreSQL с миграциями (goose) и готовой структурой для sqlc
- **Аутентификация**: JWT токены с поддержкой refresh токенов
- **Workflow**: Интеграция с Temporal для обработки асинхронных процессов
- **CLI**: Cobra для управления командами
- **Мониторинг**: Prometheus метрики и Jaeger трейсинг

## Содержание

- [Требования](#требования)
- [Быстрый старт](#быстрый-старт)
- [Архитектура проекта](#архитектура-проекта)
- [Модули приложения](#модули-приложения)
- [База данных](#база-данных)
- [Конфигурация](#конфигурация)
- [API и эндпоинты](#api-и-эндпоинты)
- [Развертывание](#развертывание)
- [Task команды](#task-команды)
- [Безопасность](#безопасность)
- [Мониторинг](#мониторинг)
- [Roadmap](#roadmap)

## Требования

- Go 1.24+
- Docker 24+ (для локальной инфраструктуры)
- Task (`go install github.com/go-task/task/v3/cmd/task@latest`) — опционально

## Быстрый старт

### 1) Поднять локальную инфраструктуру (Postgres, Jaeger)

1. Перейдите в папку `deployments/local`.
2. Создайте файл `.env` со значениями (пример):

   ```env
   POSTGRES_VERSION=16
   DATA_PATH_HOST=./data
   POSTGRES_ENTRYPOINT_INITDB=./postgres/docker-entrypoint-initdb.d
   POSTGRES_PORT=5432
   POSTGRES_DB=template
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   ```

3. Запустите инфраструктуру:

   ```bash
   docker compose up -d
   ```

Postgres: `localhost:5432`. Jaeger UI: `http://localhost:16686`.

### 2) Настроить конфиг приложения

- Скопируйте пример конфига:

  ```bash
  cp app/config/config.example.yml app/config/config.yml
  ```

- Отредактируйте `app/config/config.yml` под ваше окружение. Минимум обновите блок `PG` и `server.env`.

Пример (локально через Docker Postgres):

```yaml
server:
  host: localhost
  port: 8080
  env: local
  context_timeout: 1m
  read_header_timeout: 60s
  grace_ful_timeout: 8s

log:
  level: DEBUG

PG:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  scheme: public
  db: template

JWT:
  token_ttl: 20000m
  refresh_token_ttl: 20000m
  sign_key: secret
```

## Запуск приложения

- Запуск HTTP-сервера:

  ```bash
  go run app/main.go
  ```

  По умолчанию запускается команда `server` (cobra) и поднимается fx-граф. В режиме `local` сервер печатает зарегистрированные роуты.

- Сборка бинаря:

  ```bash
  go build -o bin/server app/main.go
  ./bin/server
  ```

## Миграции БД

Используется `goose` с встраиваемыми миграциями (`app/db/migrations/*.sql`).

- Применить миграции:

  ```bash
  go run app/main.go migration up
  ```

- Откатить миграции на 1 шаг:

  ```bash
  go run app/main.go migration down 1
  ```

Подключение к БД формируется из `app/config/config.yml` (секция `PG`). Для окружения `test` действует защита — имя БД должно содержать `test`.

## Развертывание

### Локальная разработка

1. **Инфраструктура**:
   ```bash
   cd deployments/local
   docker compose up -d
   ```

2. **Миграции**:
   ```bash
   task migrate:up
   # или
   go run app/main.go migrate up
   ```

3. **Запуск приложения**:
   ```bash
   task deps:app
   # или
   go run app/main.go
   ```

### Docker Compose сервисы

- **PostgreSQL**: `localhost:5432`
- **Jaeger UI**: `http://localhost:16686`
- **Temporal**: `localhost:7233` (опционально)

## Task команды

### Основные команды

- `task format` - Форматирование кода (gofumpt + gci)
- `task lint` - Линтинг кода (golangci-lint)
- `task ogen:gen` - Генерация кода по OpenAPI
- `task deps:update` - Обновление зависимостей

### Команды для работы с БД

- `task migrate:up` - Применить миграции
- `task migrate:down` - Откатить миграцию
- `task migrate:status` - Статус миграций
- `task migrate:create -- <name>` - Создать миграцию

### Инфраструктурные команды

- `task deps` - Запустить инфраструктуру
- `task deps:stop` - Остановить инфраструктуру
- `task deps:reset-psql` - Пересоздать PostgreSQL

## API и эндпоинты

Базовый адрес: `http://localhost:8080`

### Аутентификация

Все защищенные эндпоинты требуют Bearer токен в заголовке:
```
Authorization: Bearer <jwt_token>
```

### Основные эндпоинты

#### Логин
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "123456"
}
```

#### Информация о пользователе
```http
GET /auth/me
Authorization: Bearer <token>
```

#### Обновление токена
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "<refresh_token>"
}
```

#### Тестовый workflow
```http
POST /workflows/test-user-onboarding
Content-Type: application/json

{
  "user_id": "uuid",
  "email": "user@example.com"
}
```

### Коды ответов

- `200` - Успешный запрос
- `400` - Ошибка валидации
- `401` - Неавторизован
- `403` - Доступ запрещен
- `500` - Внутренняя ошибка сервера

OpenAPI-спека: `shared/api/app/v1/app.v1.swagger.yml`.
Сгенерированный код (server/client/middleware): `shared/pkg/openapi/app/v1/*` (генерируется `ogen`).

## Кодоген и линтеры/формат

- Генерация кода по OpenAPI (ogen):

  ```bash
  task ogen:gen
  ```

  Выход: `shared/pkg/openapi/app/v1`.

- Форматирование кода:

  ```bash
  task format
  ```

- Линтинг:

  ```bash
  task lint
  ```

- Обновление зависимостей во всех модулях:

  ```bash
  task deps:update
  ```

- Protobuf (если используется `shared/proto`):

  ```bash
  task proto:lint
  task proto:gen
  ```

## Архитектура проекта

### Структура директорий

```text
├── app/                          # Основное приложение
│   ├── cmd/                      # Команды и dependency injection
│   │   ├── dependency/           # Провайдеры для DI контейнера
│   │   └── migrateCmd/          # Команды для миграций БД
│   ├── config/                   # Конфигурация приложения
│   ├── db/migrations/           # SQL миграции
│   ├── internal/                # Внутренняя логика приложения
│   │   ├── module/              # Бизнес-модули
│   │   │   └── users/           # Модуль пользователей
│   │   ├── shared/              # Общие компоненты
│   │   │   └── middleware/      # HTTP middleware
│   │   └── workflows/           # Temporal workflows
│   └── pkg/                     # Переиспользуемые пакеты
├── shared/                      # Общие ресурсы между модулями
│   ├── api/                     # OpenAPI спецификации
│   └── pkg/openapi/            # Сгенерированный код
├── platform/                   # Платформенные утилиты
├── deployments/                 # Конфигурации для развертывания
│   └── local/                   # Локальная инфраструктура (Docker)
└── bin/                        # Скомпилированные бинарники
```

### Используемые технологии

#### Core
- **Go 1.25.1**: Основной язык программирования
- **fx**: Dependency injection контейнер
- **cobra**: CLI framework
- **chi**: HTTP роутер

#### База данных
- **PostgreSQL**: Основная СУБД
- **pgx/v5**: PostgreSQL драйвер
- **goose**: Система миграций
- **sqlc**: Готовая структура для кодогенерации SQL

#### API и кодогенерация
- **ogen**: OpenAPI кодогенерация
- **OpenAPI 3.0.3**: Спецификация API

#### Безопасность
- **JWT**: Аутентификация и авторизация
- **Argon2ID**: Хеширование паролей
- **golang-jwt/jwt/v5**: JWT library

#### Мониторинг и трейсинг
- **Prometheus**: Метрики
- **OpenTelemetry**: Трейсинг
- **Jaeger**: Distributed tracing

#### Workflow
- **Temporal**: Workflow orchestration

#### Валидация и конфигурация
- **cleanenv**: Конфигурация из ENV/YAML
- **validator/v10**: Валидация структур

## Модули приложения

### 1. Модуль пользователей (`app/internal/module/users/`)

**Назначение**: Управление пользователями, аутентификация и авторизация

**Компоненты**:
- **Controller** (`controller/v1/`): HTTP handlers совместимые с ogen
- **Service** (`services/`): Бизнес-логика работы с пользователями
- **Repository** (`repository/`): Слой доступа к данным
- **DTO** (`tokenDTO/`): Data Transfer Objects

**API Endpoints**:
- `POST /auth/login` - Авторизация пользователя
- `POST /auth/logout` - Выход из системы
- `POST /auth/refresh` - Обновление токенов
- `GET /auth/me` - Информация о текущем пользователе

### 2. Модуль Workflows (`app/internal/workflows/`)

**Назначение**: Обработка асинхронных процессов через Temporal

**Компоненты**:
- **Workflows** (`workflows/`): Определения workflow процессов
- **Activities** (`activities/`): Отдельные шаги workflow

**API Endpoints**:
- `POST /workflows/test-user-onboarding` - Тестовый запуск workflow пользователя

## База данных

### Схема БД

#### Таблица `users`
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50) DEFAULT 'user',
    tenant_id VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN DEFAULT FALSE,
    settings JSONB DEFAULT '{}',
    components JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

**Индексы**:
- `idx_users_email` - для поиска по email
- `idx_users_tenant_id` - для мультитенантности
- `idx_users_role` - для фильтрации по ролям
- `idx_users_is_active` - для фильтрации активных пользователей

**Триггеры**:
- Автоматическое обновление `updated_at` при изменении записи

### Тестовые данные

В базу автоматически добавляются тестовые пользователи:
- `admin@example.com` (пароль: `password`, роль: `admin`)
- `user@example.com` (пароль: `123456`, роль: `user`)

## Конфигурация

### Основной конфиг (`config.yml`)

```yaml
server:
  host: localhost
  port: 8080
  env: local
  context_timeout: 1m
  read_header_timeout: 60s
  grace_ful_timeout: 8s

log:
  level: DEBUG

PG:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  scheme: public
  db: template

JWT:
  token_ttl: 20000m
  refresh_token_ttl: 20000m
  sign_key: secret
```

### Переменные окружения (`.env`)

Поддерживаются переменные для:
- Конфигурации приложения (`APP_*`)
- Настроек PostgreSQL (`APP_*_PG`)
- JWT параметров (`APP_SIGN_KEY`, `APP_TOKEN_TTL`)
- Temporal настроек (`APP_TEMPORAL_*`)
- RabbitMQ (`APP_AMQP_URI`)

## Безопасность

### Аутентификация
- JWT токены с настраиваемым временем жизни
- Refresh токены для обновления сессии
- Хеширование паролей через Argon2ID

### Валидация
- Валидация входных данных через validator/v10
- Санитизация SQL запросов через pgx

### Мультитенантность
- Поддержка `tenant_id` для изоляции данных

## Мониторинг

### Метрики
- HTTP метрики (запросы, время ответа, коды ошибок)
- Метрики базы данных
- Кастомные бизнес-метрики

### Трейсинг
- OpenTelemetry интеграция
- Distributed tracing через Jaeger
- Трейсинг SQL запросов

### Логирование
- Структурированное логирование через slog
- Настраиваемый уровень логирования
- Контекстное логирование с trace ID

## Расширение проекта

### Добавление нового модуля

1. Создайте директорию в `app/internal/module/`
2. Реализуйте слои: controller, service, repository
3. Создайте fx модуль
4. Добавьте модуль в DI контейнер

### Добавление нового API эндпоинта

1. Обновите OpenAPI спецификацию в `shared/api/app/v1/`
2. Сгенерируйте код: `task ogen:gen`
3. Реализуйте handler в соответствующем контроллере
4. Добавьте business-логику в сервис

### Добавление миграции

```bash
task migrate:create -- "description_of_migration"
```

## Roadmap

### Реализовано ✅
- Logger (slog)
- CLI (cobra)
- Config (cleanenv)
- Web (chi)
- DI/IOC (fx)
- Database Postgres
- sqlc-ready структура
- Codegen (ogen)
- Migrate (goose)
- Docker Compose (локальная разработка)
- JWT аутентификация
- Temporal workflows

### Планируется 🔄
- Seed данные
- Redis интеграция
- RabbitMQ
- Docker Compose для продакшена
- Grafana дашборды
- CI/CD pipeline
- Тестовое покрытие
- Документация API

---

## Контакты и поддержка

Для вопросов и предложений по улучшению проекта создавайте issues в репозитории.
