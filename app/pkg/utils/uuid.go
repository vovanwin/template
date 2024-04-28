package utils

import "github.com/google/uuid"

func CreateUUID() uuid.UUID {
	id, _ := uuid.NewV7()
	return id
}
