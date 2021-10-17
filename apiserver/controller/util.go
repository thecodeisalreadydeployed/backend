package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
)

func writeResponse(ctx *fiber.Ctx, data interface{}, err error) error {
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}

	return ctx.JSON(data)
}
