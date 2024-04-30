package users

import (
	usersv1 "github.com/vovanwin/template/internal/module/users/controller/v1"
	"go.uber.org/fx"
)

var Module = fx.Module("authModule",
	//контроллер
	fx.Invoke(usersv1.Controller),

	//fx.Provide(
	//	//service
	//	service.NewAuthServiceImpl,
	//	// repository
	//	repository.NewEntAuthRepo,
	//),
)
