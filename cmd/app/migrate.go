package app

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"log/slog"
	"os"
	"template/config"
)

// @see https://github.com/pressly/goose/blob/master/examples/go-migrations/main.go
var testCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Управление миграциями",
	Long:  `пакет https://github.com/pressly `,
	Run: func(cmd *cobra.Command, args []string) {

		fx.New(migrationInject()).Run()
	},
}

func migrationInject() fx.Option {
	return fx.Options(
		fx.Provide(
			config.NewConfig,

			//fxslog.SetupLogger(),
		),
		fx.Invoke(migrate),
	)
}

// //go:embed migrations/*.sql
// var embedMigrations embed.FS
var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func migrate(config config.Config) {

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

	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[1]
	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}
	slog.Info("goose %v: %v", command, err)
	if err := goose.RunContext(context.Background(), command, db, "migrations", arguments...); err != nil {
		slog.Error("goose %v: %v", command, err)
	}

	slog.Info("Команда выполнена")

	os.Exit(0)
}
