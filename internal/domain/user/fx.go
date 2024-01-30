package user

import (
	"go.uber.org/fx"
	"template/internal/domain/user/controller"
	userRep "template/internal/domain/user/repository/user"
	authSer "template/internal/domain/user/service/auth"
	userSer "template/internal/domain/user/service/user"
)

var Module = fx.Module("userDomain",
	// контроллеры
	fx.Invoke(
		controller.InitIndexController,
	),
	fx.Provide(
		//репозитории
		userRep.NewPgUserRepo,

		//сервисы
		userSer.NewUserImpl,
		authSer.NewAuthImpl,
	),
)
