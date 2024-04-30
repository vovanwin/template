package migrateCmd

import (
	"context"
	"fmt"
	"github.com/vovanwin/template/config"
	"log"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/spf13/cobra"
	"os"
)

var (
	MigrationsCmd = &cobra.Command{
		Use:   "migration",
		Short: "Запуск миграций",
		Run:   migration,
	}
)

// @see https://atlasgo.io/integrations/go-sdk
func migration(cmd *cobra.Command, args []string) {

	con, _ := config.NewConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		con.PG.User, con.PG.Password, con.PG.Host, con.PG.Port, con.PG.Db, con.PG.Scheme)

	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS("./database/migrations"),
		),
	)
	if err != nil {
		log.Fatalf("failed to load working directory: %v", err)
	}
	// atlasexec works on a temporary directory, so we need to close it
	defer workdir.Close()

	// Initialize the client.
	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}
	// Применить миграции
	res, err := client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: connStr,
	})
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	fmt.Printf("Applied %d migrations\n", len(res.Applied))
}
