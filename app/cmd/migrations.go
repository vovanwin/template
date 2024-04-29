package cmd

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/database"

	"github.com/pressly/goose/v3"

	"github.com/spf13/cobra"
	_ "github.com/vovanwin/template/database"
	"log/slog"
	"os"
)

var pathMigrations = "migrations"

var (
	MigrationsCmd = &cobra.Command{
		Use:   "migration",
		Short: "Запуск миграций",
		Run:   migration,
	}
)

// @see https://github.com/pressly/goose/blob/master/examples/go-migrations/main.go
func migration(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		slog.Info("Не указаны аргументы")
		os.Exit(0)
	}

	con, _ := config.NewConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		con.PG.User, con.PG.Password, con.PG.Host, con.PG.Port, con.PG.Db, con.PG.Scheme)

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
	goose.SetBaseFS(database.EmbedMigrations)

	if err := goose.RunContext(context.Background(), args[0], db, pathMigrations, args[1:]...); err != nil {
		slog.Error("goose %v: %v", "command up: ", err)
	}

	slog.Info("Команда выполнена")

	os.Exit(0)
}
