package controller

import (
	"go.uber.org/fx"
	"template/internal/controller/auth"
)

var Module = fx.Invoke(
	auth.InitIndexController,
)
