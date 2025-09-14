package migrator

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/vovanwin/platform/pkg/logger"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// noOpLogger заглушка для goose логирования
type noOpLogger struct{}

func (n *noOpLogger) Printf(_ string, _ ...interface{}) {
	// Не делаем ничего - отключаем логирование goose
}

func (n *noOpLogger) Fatalf(_ string, _ ...interface{}) {
	// Не делаем ничего - отключаем логирование goose
}

// Migrator представляет мигратор базы данных
type Migrator struct {
	config         *Config
	db             *sql.DB
	migrationsFS   embed.FS
	migrationsPath string
	logger         interface {
		Info(ctx context.Context, msg string, fields ...zap.Field)
		Error(ctx context.Context, msg string, fields ...zap.Field)
		GetZapLogger() *zap.Logger
	}
	options *Options
}

// Options содержит дополнительные опции для мигратора
type Options struct {
	// AllowMissing разрешает пропущенные миграции
	AllowMissing bool
	// TestSafety включает проверку безопасности для тестовых БД
	TestSafety bool
	// TestKeyword ключевое слово для определения тестовой БД
	TestKeyword string
	// Timeout таймаут для операций с БД
	Timeout time.Duration
	// NoTransaction отключает транзакции для миграций
	NoTransaction bool
}

// DefaultOptions возвращает опции по умолчанию
func DefaultOptions() *Options {
	return &Options{
		AllowMissing:  true,
		TestSafety:    true,
		TestKeyword:   "test",
		Timeout:       5 * time.Minute,
		NoTransaction: false,
	}
}

// New создает новый мигратор
func New(config *Config, migrationsFS embed.FS, migrationsPath string, opts *Options) (*Migrator, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	lg := logger.Named("migrator")

	// Проверка безопасности для тестовых БД
	if opts.TestSafety && !strings.Contains(config.Database, opts.TestKeyword) {
		return nil, fmt.Errorf("database name must contain '%s' keyword for safety", opts.TestKeyword)
	}

	// Подключение к БД
	db, err := sql.Open("pgx", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Установка таймаутов
	db.SetConnMaxLifetime(opts.Timeout)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	lg.Info(ctx, "Database connection established",
		zap.String("database", config.Database),
		zap.String("schema", config.Schema))

	return &Migrator{
		config:         config,
		db:             db,
		migrationsFS:   migrationsFS,
		migrationsPath: migrationsPath,
		logger:         lg,
		options:        opts,
	}, nil
}

// Close закрывает соединение с БД
func (m *Migrator) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// Up выполняет миграции вверх
func (m *Migrator) Up(ctx context.Context) error {
	return m.runMigration(ctx, "up")
}

// Down выполняет откат миграций
func (m *Migrator) Down(ctx context.Context) error {
	return m.runMigration(ctx, "down")
}

// DownTo откатывает до определенной версии
func (m *Migrator) DownTo(ctx context.Context, version int64) error {
	return m.runMigrationWithArgs(ctx, "down-to", []string{fmt.Sprintf("%d", version)})
}

// Status показывает статус миграций
func (m *Migrator) Status(ctx context.Context) error {
	return m.runMigration(ctx, "status")
}

// Version показывает текущую версию миграций
func (m *Migrator) Version(ctx context.Context) (int64, error) {
	if err := m.setupGoose(); err != nil {
		return 0, err
	}

	version, err := goose.GetDBVersionContext(ctx, m.db)
	if err != nil {
		return 0, fmt.Errorf("failed to get database version: %w", err)
	}

	return version, nil
}

// Create создает новую миграцию
func (m *Migrator) Create(ctx context.Context, name, migrationType string) error {
	args := []string{name}
	if migrationType != "" {
		args = append(args, migrationType)
	}
	return m.runMigrationWithArgs(ctx, "create", args)
}

// Fix исправляет последовательность миграций
func (m *Migrator) Fix(ctx context.Context) error {
	return m.runMigration(ctx, "fix")
}

// Validate проверяет валидность миграций
func (m *Migrator) Validate(ctx context.Context) error {
	return m.runMigration(ctx, "validate")
}

// runMigration выполняет миграцию без дополнительных аргументов
func (m *Migrator) runMigration(ctx context.Context, command string) error {
	return m.runMigrationWithArgs(ctx, command, nil)
}

// runMigrationWithArgs выполняет миграцию с дополнительными аргументами
func (m *Migrator) runMigrationWithArgs(ctx context.Context, command string, args []string) error {
	startTime := time.Now()

	m.logger.Info(ctx, "Starting migration",
		zap.String("command", command),
		zap.Strings("args", args),
		zap.String("database", m.config.Database))

	if err := m.setupGoose(); err != nil {
		return err
	}

	// Создание контекста с таймаутом
	migrationCtx, cancel := context.WithTimeout(ctx, m.options.Timeout)
	defer cancel()

	// Настройка опций goose
	var gooseOpts []goose.OptionsFunc
	if m.options.AllowMissing {
		gooseOpts = append(gooseOpts, goose.WithAllowMissing())
	}
	// WithNoTransaction удалена из новых версий goose
	// if m.options.NoTransaction {
	//     gooseOpts = append(gooseOpts, goose.WithNoTransaction())
	// }

	// Выполнение команды
	err := goose.RunWithOptionsContext(
		migrationCtx,
		command,
		m.db,
		m.migrationsPath,
		args,
		gooseOpts...,
	)

	duration := time.Since(startTime)

	if err != nil {
		m.logger.Error(ctx, "Migration failed",
			zap.String("command", command),
			zap.Duration("duration", duration),
			zap.Error(err))
		return fmt.Errorf("migration %s failed: %w", command, err)
	}

	m.logger.Info(ctx, "Migration completed successfully",
		zap.String("command", command),
		zap.Duration("duration", duration))

	return nil
}

// setupGoose настраивает goose
func (m *Migrator) setupGoose() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	goose.SetBaseFS(m.migrationsFS)
	// Используем стандартный вывод для команд status и других
	// goose.SetLogger(&noOpLogger{}) - не отключаем, чтобы видеть результат команд

	return nil
}

// GetMigrationFiles возвращает список файлов миграций
func (m *Migrator) GetMigrationFiles() ([]string, error) {
	entries, err := m.migrationsFS.ReadDir(m.migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}
