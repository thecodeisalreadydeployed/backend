package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
)

func NewBuildScriptController(api fiber.Router) {
	api.Post("/validate", validateBuildScript)
}

func validateBuildScript(ctx *fiber.Ctx) error {
	input := dto.ValidateBuildScriptRequest{}

	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}

	if validationErrors := validator.CheckStruct(input); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}
	return ctx.JSON(map[string]string{
		"ok": cast.ToString(true),
	})
}
