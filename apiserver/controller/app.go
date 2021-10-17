package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewAppController(api fiber.Router) {
	api.Get("/", listApps)
	api.Get("/:appID", getApp)
}

func listApps(ctx *fiber.Ctx) error {
	result, err := datastore.GetAllApps()
	return writeResponse(ctx, result, err)
}

func getApp(ctx *fiber.Ctx) error {
	projectID := ctx.Params("appID")
	result, err := datastore.GetAppByID(projectID)
	return writeResponse(ctx, result, err)
}
