package usersGenv1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthLoginPost(t *testing.T) {
	// Создаем клиента
	client, err := usersGenv1.NewClient("http://localhost:8080/api/v1", nil)
	require.NoError(t, err)

	// Выполняем запрос
	ctx := context.Background()
	request := &usersGenv1.LoginRequest{Username: "123", Password: "12345678"}
	params := usersGenv1.AuthLoginPostParams{}
	token, err := client.AuthLoginPost(ctx, request, params)
	require.NoError(t, err)
	require.NotNil(t, token)

}

func TestAuthMeGet(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/auth/me", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"id": "test-id", "login": "test-user"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	// Создаем клиента
	client, err := usersGenv1.NewClient(server.URL, nil)
	require.NoError(t, err)

	// Выполняем запрос
	ctx := context.Background()
	params := usersGenv1.AuthMeGetParams{}
	user, err := client.AuthMeGet(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "test-id", user.ID)
	require.Equal(t, "test-user", user.Email)
}
