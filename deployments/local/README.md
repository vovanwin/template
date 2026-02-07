# Local infrastructure

## Подготовка

Создайте файл `.env` рядом с этим README:

```env
POSTGRES_VERSION=16
DATA_PATH_HOST=./data
POSTGRES_ENTRYPOINT_INITDB=./postgres/docker-entrypoint-initdb.d
POSTGRES_PORT=5432
POSTGRES_DB=template
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
```

## Основная инфраструктура (PostgreSQL + Temporal)

```bash
task deps          # запуск
task deps:stop     # остановка
task deps:logs     # логи
```

## Observability (Prometheus + Grafana + Tempo + OTEL Collector)

```bash
task metrics       # запуск
task metrics:stop  # остановка
task metrics:logs  # логи
```

| Сервис | Адрес |
|--------|-------|
| PostgreSQL | `localhost:5432` |
| Grafana | `http://localhost:3000` (admin/admin) |
| Prometheus | `http://localhost:9090` |
| Tempo | `http://localhost:3200` |
| OTEL Collector | `localhost:4317` (OTLP gRPC) |
