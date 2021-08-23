package apiserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/datastore"
)

func mapStatusCode(datastoreError error) int {
	switch datastoreError {
	case datastore.ErrAlreadyExists:
		return fiber.StatusConflict
	case datastore.ErrInvalidArgument:
		return fiber.StatusUnprocessableEntity
	case datastore.ErrFailedPrecondition:
		return fiber.StatusBadRequest
	case datastore.ErrNotFound:
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}
