package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

func NewBuildScriptController(api fiber.Router) {
	api.Post("/validate", validateBuildScript)
}

func validateBuildScript(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{
		"ok": cast.ToString(true),
	})
}
