package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewHealthController(api fiber.Router) {
	api.Get("/", getHealth)
}

func getHealth(ctx *fiber.Ctx) error {
	ok := datastore.IsReady()
	return ctx.JSON(map[string]string{
		"ok": cast.ToString(ok),
	})
}
