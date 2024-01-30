package seed

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreate, downCreate)
	goose.WithNoVersioning()

}

func upCreate(ctx context.Context, tx *sql.Tx) error {

	users := []struct {
		username string
		password string
	}{
		{
			username: "testuser1",
			password: "testuser1",
		},
		{
			username: "testuser2",
			password: "testuser2",
		},
		{
			username: "testuser2",
			password: "testuser2",
		},
	}

	for _, user := range users {
		query := "INSERT INTO users (username, password) VALUES ($1, $2)"
		if _, err := tx.ExecContext(ctx, query, user.username, user.password); err != nil {
			return err
		}
	}

	return nil
}

func downCreate(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
