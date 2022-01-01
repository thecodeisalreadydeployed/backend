package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
	"strconv"
)

func NewAppController(api fiber.Router) {
	api.Get("/list", listApps)
	api.Get("/:appID", getApp)
	api.Get("/name/:appName", searchApp)
	api.Post("/:appID/deployments", createDeployment)
	api.Get("/:appID/deployments", listAppDeployments)
	api.Post("/", createApp)
	api.Delete("/:appID", deleteApp)
	api.Put("/:appID/observable/enable", enableObservable)
	api.Put("/:appID/observable/disable", disableObservable)
}

func listApps(ctx *fiber.Ctx) error {
	result, err := datastore.GetAllApps(datastore.GetDB())
	return writeResponse(ctx, result, err)
}

func getApp(ctx *fiber.Ctx) error {
	appID := ctx.Params("appID")
	result, err := datastore.GetAppByID(datastore.GetDB(), appID)
	return writeResponse(ctx, result, err)
}

func searchApp(ctx *fiber.Ctx) error {
	appName := ctx.Params("appName")
	result, err := datastore.GetAppsByName(datastore.GetDB(), appName)
	return writeResponse(ctx, result, err)
}

func listAppDeployments(ctx *fiber.Ctx) error {
	appID := ctx.Params("appID")
	result, err := datastore.GetDeploymentsByAppID(datastore.GetDB(), appID)
	return writeResponse(ctx, result, err)

}

func createApp(ctx *fiber.Ctx) error {
	input := dto.CreateAppRequest{}

	if err := validator.ParseBodyAndValidate(ctx, &input); err != nil {
		return err
	}

	inputModel := input.ToModel()
	app, createErr := datastore.SaveApp(datastore.GetDB(), &inputModel)

	return writeResponse(ctx, app, createErr)
}

func deleteApp(ctx *fiber.Ctx) error {
	appID := ctx.Params("appID")
	err := datastore.RemoveApp(datastore.GetDB(), appID)
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func enableObservable(ctx *fiber.Ctx) error {
	appID := ctx.Params("appID")
	err := datastore.SetObservable(datastore.GetDB(), appID, true)
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func disableObservable(ctx *fiber.Ctx) error {
	appID := ctx.Params("appID")
	err := datastore.SetObservable(datastore.GetDB(), appID, false)
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.SendStatus(fiber.StatusOK)
}
