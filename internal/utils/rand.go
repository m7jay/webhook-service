package utils

import (
	"github.com/google/uuid"
)

func GenerateUniqueID() string {
	return uuid.New().String()
}
