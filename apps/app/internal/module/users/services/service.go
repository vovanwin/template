package service

//go:generate mockgen -source=$GOFILE -destination=mocks/service_mock.gen.go -package=usersServiceMocks

var _ UsersService = (*UsersServiceImpl)(nil)

type (
	UsersService interface {
	}
	UsersServiceImpl struct {
	}
)

func NewUsersServiceImpl() UsersService {
	return &UsersServiceImpl{}
}
