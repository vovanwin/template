package cmd

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"template/config"
)

// @see https://github.com/pressly/goose/blob/master/examples/go-migrations/main.go
var testCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Управление миграциями",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := config.NewConfig
		migrate(config(), args)
	},
}

func migrate(config config.Config, args []string) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.Scheme,
	)

	var db *sql.DB
	// setup database
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		slog.Error("Unable to connect to database because %s", err)
	}

	if err = db.Ping(); err != nil {
		slog.Error("Cannot ping database because %s", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.RunContext(context.Background(), args[0], db, "database/migrations", args[1:]...); err != nil {
		slog.Error("goose %v: %v", "command up: ", err)
	}

	slog.Info("Команда выполнена")

	os.Exit(0)
}
