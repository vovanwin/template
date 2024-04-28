package utils

import (
	"errors"
	"fmt"
)

var (
	// Для репозиториев
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")

	// Для сервисов
	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("не валидный токен")
	ErrValidPassword    = fmt.Errorf("cannot parse token")
	ErrRefreshToken     = fmt.Errorf("ошибка создания рефреш токена")

	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrCannotGetUser     = fmt.Errorf("cannot get user")

	ErrAccountAlreadyExists = fmt.Errorf("account already exists")
	ErrCannotCreateAccount  = fmt.Errorf("cannot create account")
	ErrAccountNotFound      = fmt.Errorf("account not found")
	ErrCannotGetAccount     = fmt.Errorf("cannot get account")

	ErrCannotCreateReservation = fmt.Errorf("cannot create reservation")

	//reports
	ErrFormReports = fmt.Errorf("отчет в процессе формирования")

	//policy
	ErrForbidden = fmt.Errorf("forbidden")
)
