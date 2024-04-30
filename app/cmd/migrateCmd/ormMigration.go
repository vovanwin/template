package migrateCmd

import (
	atlas "ariga.io/atlas/sql/migrate"
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"

	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/shared/store/gen/migrate"
	"log"
)

var (
	CreateMigrationCmd = &cobra.Command{
		Use:   "migration:orm",
		Short: "Создать миграции из моделей ORM Ent",
		Run:   createMigration,
	}
)

// @see https://entgo.io/docs/versioned-migrations
func createMigration(cmd *cobra.Command, args []string) {

	ctx := context.Background()
	// Создайте локальный каталог миграции, способный понимать формат файла миграции goose для воспроизведения.
	dir, err := atlas.NewLocalDir("database/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}
	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
	}
	if len(args) != 1 {
		log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
	}

	con, _ := config.NewConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		con.PG.User, con.PG.Password, con.PG.Host, con.PG.Port, con.PG.Db, "test") // тут указываем тустовую Бд чтобы проверить разницу между sql и текущей схемой, так устроена работа atlas

	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	err = migrate.NamedDiff(ctx, connStr, args[0], opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
