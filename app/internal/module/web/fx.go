package web

import "go.uber.org/fx"

var Module = fx.Module(
	"webModule",
	fx.Invoke(Controller),
)
