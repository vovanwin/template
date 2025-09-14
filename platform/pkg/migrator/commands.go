package migrator

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// NewMigrateCommand создает команду для миграций
func NewMigrateCommand(migrator *Migrator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
		Long:  "Execute database migrations using goose",
	}

	// Добавляем подкоманды
	cmd.AddCommand(
		newUpCommand(migrator),
		newDownCommand(migrator),
		newDownToCommand(migrator),
		newStatusCommand(migrator),
		newVersionCommand(migrator),
		newCreateCommand(migrator),
		newFixCommand(migrator),
		newValidateCommand(migrator),
	)

	return cmd
}

// newUpCommand создает команду up
func newUpCommand(migrator *Migrator) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Migrate the database to the most recent version available",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			return migrator.Up(ctx)
		},
	}
}

// newDownCommand создает команду down
func newDownCommand(migrator *Migrator) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Roll back the version by 1",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			return migrator.Down(ctx)
		},
	}
}

// newDownToCommand создает команду down-to
func newDownToCommand(migrator *Migrator) *cobra.Command {
	return &cobra.Command{
		Use:   "down-to [version]",
		Short: "Roll back to a specific version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid version number: %w", err)
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			return migrator.DownTo(ctx, version)
		},
	}
}

// newStatusCommand создает команду status
func newStatusCommand(migrator *Migrator) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Print the status of all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			return migrator.Status(ctx)
		},
	}
}

// newVersionCommand создает команду version
func newVersionCommand(migrator *Migrator) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the current version of the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			version, err := migrator.Version(ctx)
			if err != nil {
				return err
			}

			fmt.Printf("Current database version: %d\n", version)
			return nil
		},
	}
}

// newCreateCommand создает команду create
func newCreateCommand(migrator *Migrator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			migrationType, _ := cmd.Flags().GetString("type")

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			return migrator.Create(ctx, args[0], migrationType)
		},
	}

	cmd.Flags().StringP("type", "t", "sql", "Migration type (sql or go)")
	return cmd
}

// newFixCommand создает команду fix
func newFixCommand(migrator *Migrator) *cobra.Command {
	return &cobra.Command{
		Use:   "fix",
		Short: "Apply sequential ordering to migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			return migrator.Fix(ctx)
		},
	}
}

// newValidateCommand создает команду validate
func newValidateCommand(migrator *Migrator) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate migration files without running them",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			return migrator.Validate(ctx)
		},
	}
}
