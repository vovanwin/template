package controller

import (
	"go.uber.org/fx"
	"template/internal/controller/user"
)

var Module = fx.Invoke(
	user.InitIndexController,
)
