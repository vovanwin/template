package users

import (
	usersv1 "github.com/vovanwin/template/app/internal/module/users/controller/v1"
	"github.com/vovanwin/template/app/internal/module/users/repository"
	service "github.com/vovanwin/template/app/internal/module/users/services"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"authModule",
	//контроллер
	fx.Invoke(usersv1.Controller),

	fx.Provide(
		//repository
		repository.NewPostgresUsersRepository,
		//service
		service.NewUsersServiceImpl,
	),
)
