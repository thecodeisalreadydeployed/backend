package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewHealthController(api fiber.Router, dataStore datastore.DataStore) {
	api.Get("/", getHealth(dataStore))
}

func getHealth(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ok := dataStore.IsReady()
		return ctx.JSON(map[string]string{
			"ok": cast.ToString(ok),
		})
	}
}
