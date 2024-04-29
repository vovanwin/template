package auth

import (
	controller "github.com/vovanwin/template/internal/module/auth/controllers"
	"github.com/vovanwin/template/internal/module/auth/repository"
	service "github.com/vovanwin/template/internal/module/auth/services"
	"go.uber.org/fx"
)

var Module = fx.Module("authModule",
	//контроллер
	fx.Invoke(controller.Controller),

	fx.Provide(
		//service
		service.NewAuthServiceImpl,
		// repository
		repository.NewEntAuthRepo,
	),
)
