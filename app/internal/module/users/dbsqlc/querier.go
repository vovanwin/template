// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package dbsqlc

import (
	"context"

	"app/internal/shared/types"
)

type Querier interface {
	// Получить пользователя для me запроса
	FindMeForId(ctx context.Context, id types.UserID) (*Users, error)
}

var _ Querier = (*Queries)(nil)
