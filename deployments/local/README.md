# Local infrastructure

Локальная инфраструктура для разработки: Postgres и Jaeger.

## Подготовка

Создайте файл `.env` рядом с этим README со значениями (пример):

```env
POSTGRES_VERSION=16
DATA_PATH_HOST=./data
POSTGRES_ENTRYPOINT_INITDB=./postgres/docker-entrypoint-initdb.d
POSTGRES_PORT=5432
POSTGRES_DB=template
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
```

## Запуск

```bash
docker compose up -d
```

- Postgres: `localhost:5432`
- Jaeger UI: `http://localhost:16686`

## Остановка

```bash
docker compose down -v
```
