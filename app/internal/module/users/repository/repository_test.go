package repository_test

import (
	"context"
	"github.com/vovanwin/template/internal/shared/types"
	"log"
	"os"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/vovanwin/template/internal/module/users/repository"
	"github.com/vovanwin/template/internal/module/users/tokenDTO"
	"github.com/vovanwin/template/internal/shared/store/gen"
	"github.com/vovanwin/template/pkg/framework"
	"github.com/vovanwin/template/pkg/utils"
)

// Setup the test database
func setupTestDB() (*gen.Client, func()) {
	drv, err := sql.Open(dialect.SQLite, "file:test.db?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	client := gen.NewClient(gen.Driver(drv))

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client, func() {
		client.Close()
		os.Remove("file:test.db?mode=memory&cache=shared&_fk=1")
	}
}

func TestEntUsersRepo_GetMe(t *testing.T) {
	client, cleanup := setupTestDB()
	defer cleanup()

	repo := repository.NewEntUsersRepo(client)

	t.Run("successful get me", func(t *testing.T) {
		userID := types.NewUserID()
		expectedUser := &gen.Users{
			ID:        userID,
			Login:     "testuser",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		client.Users.Create().
			SetID(userID).
			SetLogin("testuser").
			SetPassword("password").
			SaveX(context.Background())

		ctx := context.WithValue(context.Background(), framework.Claims, &tokenDTO.TokenClaims{UserId: userID})

		user, err := repo.GetMe(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Login, user.Login)
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), framework.Claims, &tokenDTO.TokenClaims{UserId: types.NewUserID()})

		user, err := repo.GetMe(ctx)
		assert.Error(t, err)
		assert.Equal(t, utils.ErrNotFound, err)
		assert.Nil(t, user)
	})

	t.Run("other error", func(t *testing.T) {
		// Simulate other errors by closing the client and running the test
		client.Close()

		ctx := context.WithValue(context.Background(), framework.Claims, &tokenDTO.TokenClaims{UserId: types.NewUserID()})

		user, err := repo.GetMe(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository: get user me:")
		assert.Nil(t, user)
	})
}

func TestEntUsersRepo_FindForLogin(t *testing.T) {
	client, cleanup := setupTestDB()
	defer cleanup()

	repo := repository.NewEntUsersRepo(client)

	t.Run("successful find for login", func(t *testing.T) {
		userID := types.NewUserID()
		expectedUser := &gen.Users{
			ID:        userID,
			Login:     "testuser",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		client.Users.Create().
			SetID(userID).
			SetLogin("testuser").
			SetPassword("password").
			SaveX(context.Background())

		user, err := repo.FindForLogin(context.Background(), "testuser")
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Login, user.Login)
	})

	t.Run("user not found", func(t *testing.T) {
		user, err := repo.FindForLogin(context.Background(), "nonexistent")
		assert.Error(t, err)
		assert.Equal(t, utils.ErrNotFound, err)
		assert.Nil(t, user)
	})

	t.Run("other error", func(t *testing.T) {
		// Simulate other errors by closing the client and running the test
		client.Close()

		user, err := repo.FindForLogin(context.Background(), "testuser")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository: get user me:")
		assert.Nil(t, user)
	})
}
