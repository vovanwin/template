package repository

import (
	"go.uber.org/fx"
	"template/internal/repository/user"
)

var Module = fx.Provide(
	user.NewPgUserRepo,
)
