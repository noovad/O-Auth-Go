package helper

import (
	"errors"

	"github.com/google/uuid"
)

func StringToUUID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, errors.New("empty ID string")
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID format")
	}

	return parsedUUID, nil
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
