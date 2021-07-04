package utils

import (
	"github.com/google/uuid"
)

type UUID string

func GenerateNewID() UUID {
	return UUID(uuid.NewString())
}