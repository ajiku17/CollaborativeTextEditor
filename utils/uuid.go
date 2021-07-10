package utils

import (
	"github.com/google/uuid"
)

type UUID string

func GenerateNewUUID() UUID {
	return UUID(uuid.NewString())
}