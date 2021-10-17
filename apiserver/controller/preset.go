package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/preset"
)

func NewPresetController(api fiber.Router) {
	api.Get("/:preset", getPreset)
}

func getPreset(ctx *fiber.Ctx) error {
	framework := ctx.Params("preset")
	text := preset.Text(preset.Framework(framework))
	return ctx.SendString(text)
}
