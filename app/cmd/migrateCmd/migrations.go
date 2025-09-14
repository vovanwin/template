package migrateCmd

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
	embeded "github.com/vovanwin/template/app/db"

	"github.com/vovanwin/template/app/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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
func migration(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Ошибка: не указан параметр команды")
		os.Exit(0)
	}

	if args[0] != "up" && args[0] != "down" {
		pathMigrations = "db/" + pathMigrations
	}

	var db *sql.DB
	con, _ := config.NewConfig()
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable&search_path=%s",
		con.PG.UserPG, con.PG.PasswordPG,
		net.JoinHostPort(con.PG.HostPG, con.PG.PortPG),
		con.PG.DbNamePG, con.PG.SchemePG,
	)

	if con.IsTest() {
		if !containsTest(con.PG.DbNamePG) {
			panic("Возможно указана не тестовая БД, Тестовая БД должна иметь в названии 'test' ")
		}
	}

	// setup database
	db, err := sql.Open("pgx", connStr)
	defer db.Close()

	if err != nil {
		fmt.Printf("Ошибка подключения к БД: %v\n", err)
	}

	if err = db.Ping(); err != nil {
		fmt.Printf("Ошибка ping БД: %v\n", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	goose.SetBaseFS(embeded.EmbedMigrations)

	if err := goose.RunWithOptionsContext(context.Background(), args[0], db, pathMigrations, args[1:], goose.WithAllowMissing()); err != nil {
		fmt.Printf("Ошибка выполнения goose: %v\n", err)
	}

	fmt.Println("Команда выполнена")

	os.Exit(0)
}

// Функция для проверки наличия слова "test" в строке, доп проверка, предполагается в название Бд для тестов есть
// слово test, так же исключит вариант когда неправильно выбрали БД и запустили тесты,
func containsTest(s string) bool {
	return strings.Contains(s, "test")
}
