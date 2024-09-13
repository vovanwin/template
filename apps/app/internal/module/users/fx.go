package users

import (
	usersv1 "app/internal/module/users/controller/v1"
	service "app/internal/module/users/services"
	"go.uber.org/fx"
)

var Module = fx.Module("authModule",
	//контроллер
	fx.Invoke(usersv1.Controller),

	fx.Provide(
		//service
		service.NewUsersServiceImpl,
	),
)
