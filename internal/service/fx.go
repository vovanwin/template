package service

import (
	"go.uber.org/fx"
	"template/internal/service/user"
)

var Module = fx.Provide(
	user.NewUserImpl,
)
