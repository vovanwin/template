# Миграции и сессии - Руководство по настройке

## 📋 Что добавлено

### 1. Миграции базы данных
- **users** - таблица пользователей с argon2id хешированием паролей  
- **sessions** - таблица для хранения сессий в PostgreSQL

### 2. Сессии через PostgreSQL  
- Настроена через `scs` с `postgresstore`
- Автоматическая очистка истекших сессий
- Возможность переключения на Redis

### 3. Пользователи и аутентификация
- Repository pattern для работы с пользователями
- Service layer с валидацией учетных данных  
- Web-контроллеры используют настоящую аутентификацию

## 🚀 Запуск

### 1. Запустите PostgreSQL
```bash
cd deployments/local
docker-compose up -d postgres
```

### 2. Выполните миграции
```bash
cd app
go run . migration up
```

### 3. Запустите приложение
```bash
go run .
```

## 🔐 Тестовые аккаунты

После выполнения миграций доступны:

- **admin@example.com** / `password`
- **user@example.com** / `123456`

## ⚙️ Конфигурация сессий

В `config/config.yml`:

```yaml
sessions:
  store: postgres      # postgres, redis, memory
  lifetime: 24h        # время жизни сессий
  # redis_addr: localhost:6379    # для Redis
  # redis_password: ""
  # redis_db: 0
```

## 🔄 Переключение на Redis

1. Добавьте в `go.mod`:
```bash
go get github.com/alexedwards/scs/redisstore
```

2. Обновите `pkg/sessions/session_provider.go`:
```go
case Redis:
    // Реализуйте Redis store
    redisPool := &redis.Pool{...}
    store := redisstore.New(redisPool)
    sessionManager.Store = store
```

3. Измените в конфиге:
```yaml
sessions:
  store: redis
  redis_addr: localhost:6379
```

## 📁 Структура файлов

```
app/
├── db/migrations/
│   ├── 20241225000001_create_users_table.sql
│   └── 20241225000002_create_sessions_table.sql
├── internal/module/users/repository/
│   └── users_repository.go
├── pkg/sessions/
│   └── session_provider.go
└── config/config.yml (обновлен)
```

## 🔍 Функциональность

- ✅ Сессии хранятся в PostgreSQL
- ✅ Автоочистка истекших сессий
- ✅ Безопасное хеширование паролей (argon2id)
- ✅ Repository pattern для пользователей
- ✅ Web-аутентификация через HTMX
- ✅ Конфигурируемый выбор хранилища сессий

## 🛠 Дополнительные команды

### Откат миграций
```bash
go run . migration down
```

### Просмотр сессий в БД
```sql
SELECT token, expiry FROM sessions;
```

### Очистка истекших сессий вручную
```sql
SELECT cleanup_expired_sessions();
```
