package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewAppController(api fiber.Router) {
	api.Get("/", listApps)
	api.Get("/:appID", getApp)
	api.Get("/:appID/deployments", getApp)
}

func listApps(ctx *fiber.Ctx) error {
	result, err := datastore.GetAllApps()
	return writeResponse(ctx, result, err)
}

func getApp(ctx *fiber.Ctx) error {
	appID := ctx.Params("appID")
	result, err := datastore.GetAppByID(appID)
	return writeResponse(ctx, result, err)
}

func listAppDeployments(ctx *fiber.Ctx) error {
	appID := ctx.Params("appID")
	result, err := datastore.GetDeploymentsByAppID(appID)
	return writeResponse(ctx, result, err)
}

func createApp(ctx *fiber.Ctx) error {
	input := dto.CreateAppRequest{}

	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}

	if validationErrors := validator.CheckStruct(input); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	inputModel := input.ToModel()
	app, createErr := datastore.SaveApp(&inputModel)

	return writeResponse(ctx, app, createErr)
}
