package errutil

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/errutil"
)

func MapStatusCode(err error) int {
	if errutil.IsAlreadyExists(err) {
		return fiber.StatusBadRequest
	}

	if errutil.IsFailedPrecondition(err) {
		return fiber.StatusPreconditionFailed
	}

	if errutil.IsInvalidArgument(err) {
		return fiber.StatusBadRequest
	}

	if errutil.IsNotFound(err) {
		return fiber.StatusNotFound
	}

	if errutil.IsNotImplemented(err) {
		return fiber.StatusNotImplemented
	}

	return fiber.StatusInternalServerError
}
