package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitapi"
	"github.com/thecodeisalreadydeployed/repositoryobserver"
	"github.com/thecodeisalreadydeployed/statusapi"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"
)

func NewAppController(
	api fiber.Router,
	workloadController workloadcontroller.WorkloadController,
	observer repositoryobserver.RepositoryObserver,
	statusAPIBackend statusapi.StatusAPIBackend,
	gitAPIBackend gitapi.GitAPIBackend,
) {
	api.Get("/list", listApps)
	api.Get("/search", searchApp)
	api.Get("/:appID", getApp)
	api.Get("/:appID/status", getAppStatus(statusAPIBackend))
	api.Post("/:appID/deployments", createDeployment(workloadController))
	api.Get("/:appID/deployments", listAppDeployments)
	api.Post("/", createApp(workloadController, gitAPIBackend))
	api.Delete("/:appID", deleteApp)
	api.Post("/:appID/observable/enable", enableObservable)
	api.Post("/:appID/observable/disable", disableObservable)
	api.Post("/:appID/refresh", forceRefresh(observer))
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

func getAppStatus(statusAPIBackend statusapi.StatusAPIBackend) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		result, err := statusAPIBackend.GetActiveDeploymentID(appID)
		return writeResponse(ctx, map[string]string{"deploymentID": result}, err)
	}
}

func searchApp(ctx *fiber.Ctx) error {
	appName := ctx.Query("name")
	result, err := datastore.GetAppsByName(datastore.GetDB(), appName)
	return writeResponse(ctx, result, err)
}

func listAppDeployments(ctx *fiber.Ctx) error {
	appID := ctx.Params("appID")
	result, err := datastore.GetDeploymentsByAppID(datastore.GetDB(), appID)
	return writeResponse(ctx, result, err)

}

func createApp(workloadController workloadcontroller.WorkloadController, gitAPIBackend gitapi.GitAPIBackend) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := dto.CreateAppRequest{}
		if err := validator.ParseBodyAndValidate(c, &input); err != nil {
			return err
		}
		inputModel := input.ToModel()
		gs, err := gitAPIBackend.FillGitSource(&(inputModel.GitSource))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		inputModel.GitSource = *gs
		app, createErr := workloadController.NewApp(&inputModel)
		return writeResponse(c, app, createErr)
	}
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

func forceRefresh(observer repositoryobserver.RepositoryObserver) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		ok := observer.Refresh(appID)
		if ok {
			return ctx.SendStatus(fiber.StatusOK)
		} else {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
	}
}
