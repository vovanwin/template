package user

import (
	"go.uber.org/fx"
	userRep "template/internal/domain/user/repository/user"
	authSer "template/internal/domain/user/service/auth"
	userSer "template/internal/domain/user/service/user"
)

var Module = fx.Module("userDomain",
	fx.Provide(
		//репозитории
		userRep.NewPgUserRepo,

		//сервисы
		userSer.NewUserImpl,
		authSer.NewAuthImpl,
	),
)
