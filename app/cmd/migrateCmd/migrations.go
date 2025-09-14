package migrateCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vovanwin/platform/pkg/migrator"
	"github.com/vovanwin/template/app/config"
	embeded "github.com/vovanwin/template/app/db"
)

var pathMigrations = "migrations"

// NewMigrateCommand создает команду migrate с использованием платформенного мигратора
func NewMigrateCommand() *cobra.Command {
	// Загружаем конфигурацию
	con, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("ошибка загрузки конфигурации: %w", err))
	}

	// Настройка опций мигратора
	opts := migrator.DefaultOptions()
	opts.TestSafety = con.IsTest()

	// Создаем адаптер конфигурации
	configAdapter := NewConfigAdapter(con)
	migratorConfig := migrator.NewConfig(configAdapter)

	// Создаем мигратор
	m, err := migrator.New(migratorConfig, embeded.EmbedMigrations, pathMigrations, opts)
	if err != nil {
		panic(fmt.Errorf("ошибка создания мигратора: %w", err))
	}

	// Возвращаем готовую команду с подкомандами
	return migrator.NewMigrateCommand(m)
}

// MigrationsCmd - alias для обратной совместимости (можно удалить позже)
var MigrationsCmd = NewMigrateCommand()
