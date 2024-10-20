package utils

import (
	"errors"
	"fmt"
)

var (
	// Для репозиториев
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	//policy.
	ErrForbidden    = fmt.Errorf("forbidden")
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrValidation   = fmt.Errorf("validation")
)
