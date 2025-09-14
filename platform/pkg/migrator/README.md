# Платформенный мигратор

Универсальная библиотека для работы с миграциями базы данных, основанная на goose.

## Особенности

- ✅ Автоматические проверки безопасности для тестовых БД
- ✅ Подробное логирование всех операций
- ✅ Поддержка таймаутов и настройки соединения
- ✅ Валидация миграций
- ✅ Создание новых миграций
- ✅ Встроенные команды Cobra
- ✅ Гибкая конфигурация

## Использование

### Базовое использование

```go
package main

import (
    "context"
    "embed"
    "log"

    "github.com/vovanwin/platform/pkg/migrator"
    "github.com/vovanwin/template/app/config"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
    // Загрузка конфигурации
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatal(err)
    }

    // Создание адаптера конфигурации
    dbConfig := &DatabaseConfigImpl{cfg: cfg}
    migratorConfig := migrator.NewConfig(dbConfig)

    // Создание мигратора
    m, err := migrator.New(migratorConfig, migrationsFS, "migrations", nil)
    if err != nil {
        log.Fatal(err)
    }
    defer m.Close()

    // Выполнение миграции
    ctx := context.Background()
    if err := m.Up(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### Использование с Cobra CLI

```go
package main

import (
    "github.com/spf13/cobra"
    "github.com/vovanwin/platform/pkg/migrator"
)

func main() {
    // Создание мигратора (как выше)...

    // Создание команд CLI
    rootCmd := &cobra.Command{Use: "myapp"}
    migrateCmd := migrator.NewMigrateCommand(m)
    rootCmd.AddCommand(migrateCmd)

    rootCmd.Execute()
}
```

## Доступные команды

### Task команды

```bash
# Выполнить миграции для сервиса app (по умолчанию)
task migrate:up

# Выполнить миграции для конкретного сервиса
SERVICE=users task migrate:up

# Откатить миграцию
task migrate:down
SERVICE=users task migrate:down

# Статус миграций
task migrate:status
SERVICE=users task migrate:status

# Текущая версия БД
task migrate:version

# Создать новую миграцию
task migrate:create -- create_users_table

# Выполнить миграции для всех сервисов
task migrate:all:up

# Статус миграций для всех сервисов
task migrate:all:status

# Версии БД для всех сервисов
task migrate:all:version
```

### CLI команды приложения

```bash
# Базовые команды
go run app/main.go migrate up
go run app/main.go migrate down
go run app/main.go migrate status
go run app/main.go migrate version

# Создание миграций
go run app/main.go migrate create create_users_table
go run app/main.go migrate create --type sql create_index

# Откат до версии
go run app/main.go migrate down-to 20241225000001

# Утилиты
go run app/main.go migrate fix      # Исправить последовательность
go run app/main.go migrate validate # Проверить файлы миграций
```

## Конфигурация

### Переменные Taskfile

Добавьте ваши сервисы в переменную `SERVICES` в `Taskfile.yaml`:

```yaml
vars:
  SERVICES: app users orders notifications
```

### Опции мигратора

```go
opts := &migrator.Options{
    AllowMissing:  true,           // Разрешить пропущенные миграции
    TestSafety:    true,           // Проверка безопасности для тестов
    TestKeyword:   "test",         // Ключевое слово для тестовой БД
    Timeout:       5*time.Minute,  // Таймаут операций
    NoTransaction: false,          // Отключить транзакции
}
```

### Адаптер конфигурации

Реализуйте интерфейс `DatabaseConfig`:

```go
type DatabaseConfigImpl struct {
    cfg *config.Config
}

func (c *DatabaseConfigImpl) GetHost() string     { return c.cfg.PG.HostPG }
func (c *DatabaseConfigImpl) GetPort() string     { return c.cfg.PG.PortPG }
func (c *DatabaseConfigImpl) GetUsername() string { return c.cfg.PG.UserPG }
func (c *DatabaseConfigImpl) GetPassword() string { return c.cfg.PG.PasswordPG }
func (c *DatabaseConfigImpl) GetDatabase() string { return c.cfg.PG.DbNamePG }
func (c *DatabaseConfigImpl) GetSchema() string   { return c.cfg.PG.SchemePG }
func (c *DatabaseConfigImpl) GetSSLMode() string  { return "disable" }
```

## Структура файлов

```
app/
├── db/
│   └── migrations/
│       └── 20241225000001_create_users_table.sql
├── cmd/
│   └── migrateCmd/
│       ├── migrations.go      # Основная команда
│       └── config_adapter.go  # Адаптер конфигурации
└── main.go

platform/
└── pkg/
    └── migrator/
        ├── config.go          # Конфигурация
        ├── migrator.go        # Основная логика
        ├── commands.go        # Cobra команды
        ├── utils.go           # Утилиты
        └── README.md
```

## Безопасность

- Автоматическая проверка имени БД на наличие слова "test" в тестовом режиме
- Валидация имен миграций
- Проверка подключения перед выполнением операций
- Контролируемые таймауты для всех операций

## Логирование

Все операции логируются с детальной информацией:
- Время начала и завершения операций
- Ошибки с полным контекстом
- Информация о подключении к БД
- Статус выполнения команд