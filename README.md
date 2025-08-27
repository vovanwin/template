# Шаблон Golang приложения

Минимальный каркас для старта проекта на Go: HTTP (chi), DI (fx), CLI (cobra), миграции (goose), кодоген по OpenAPI (ogen), Postgres (готовая структура для sqlc).

## Содержание

- Требования
- Быстрый старт (локально и Docker)
- Конфигурация
- Запуск приложения
- Миграции БД
- API и эндпоинты
- Кодоген (OpenAPI/ogen) и линтеры/форматирование
- Структура проекта
- Roadmap

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

## API и эндпоинты

Базовый адрес: `http://localhost:8080`

- Healthcheck: `GET /api/v1/healthcheck/` → `OK`
- Версия билда: `GET /version`
- Пользователь:
  - `POST /auth/login` → `AuthToken` (демо-заглушка)
  - `GET /auth/me` → `UserMe` (требует Bearer-токен; демо-ответ)

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

## Структура проекта

```text
app/
  cmd/
    dependency/        # Провайдеры: конфиг, логгер, сервер, БД
    migrateCmd/        # Команда миграций (goose)
  config/              # Конфигурация приложения
  db/
    migrations/        # SQL-миграции (встраиваемые)
  internal/
    module/
      healthcheck/     # Healthcheck и /version
      users/           # Пользовательские хендлеры/сервисы
        controller/v1/ # Ogen-совместимые хендлеры и security
        services/
    shared/
      middleware/      # Метрики, трейсинг и пр.
      validator/
  pkg/
    httpserver/        # Обертка над chi/http.Server
    response/          # Ответы/ошибки
    storage/postgres/  # Подключение к Postgres (pgx)
shared/
  api/app/v1/          # OpenAPI-спека
  pkg/openapi/app/v1/  # Сгенерированный код (ogen)

deployments/
  local/               # docker-compose для локалки (Postgres, Jaeger)
```

## Roadmap

- [x] Logger (slog)
- [x] CLI (cobra)
- [x] Config (cleanenv)
- [x] Web (chi)
- [x] DI/IOC (fx)
- [x] Database Postgres
- [x] sqlc-ready структура
- [x] Codegen (ogen)
- [x] Migrate (goose)
- [x] Docker Compose (локальная разработка)
- [ ] Seed
- [ ] Redis
- [ ] Temporal
- [ ] RabbitMQ
- [ ] Docker Compose для проды
   
