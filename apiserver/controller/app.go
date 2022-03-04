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
	dataStore datastore.DataStore,
) {
	api.Get("/list", listApps(dataStore))
	api.Get("/search", searchApp(dataStore))
	api.Get("/:appID", getApp(dataStore))
	api.Get("/:appID/status", getAppStatus(statusAPIBackend, dataStore))
	api.Post("/:appID/deployments", createDeployment(workloadController, dataStore))
	api.Get("/:appID/deployments", listAppDeployments(dataStore))
	api.Post("/", createApp(workloadController, gitAPIBackend, dataStore))
	api.Delete("/:appID", deleteApp(dataStore))
	api.Post("/:appID/observable/enable", enableObservable(dataStore))
	api.Post("/:appID/observable/disable", disableObservable(dataStore))
	api.Post("/:appID/refresh", forceRefresh(observer))
}

func listApps(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		result, err := dataStore.GetAllApps()
		return writeResponse(ctx, result, err)
	}
}

func getApp(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		result, err := dataStore.GetAppByID(appID)
		return writeResponse(ctx, result, err)
	}
}

func getAppStatus(statusAPIBackend statusapi.StatusAPIBackend, dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		result, err := statusAPIBackend.GetActiveDeploymentID(appID, dataStore)
		return writeResponse(ctx, map[string]string{"deploymentID": result}, err)
	}
}

func searchApp(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appName := ctx.Query("name")
		result, err := dataStore.GetAppsByName(appName)
		return writeResponse(ctx, result, err)
	}
}

func listAppDeployments(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		result, err := dataStore.GetDeploymentsByAppID(appID)
		return writeResponse(ctx, result, err)
	}

}

func createApp(workloadController workloadcontroller.WorkloadController, gitAPIBackend gitapi.GitAPIBackend, dataStore datastore.DataStore) fiber.Handler {
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
		app, createErr := workloadController.NewApp(&inputModel, dataStore)
		return writeResponse(c, app, createErr)
	}
}

func deleteApp(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		err := dataStore.RemoveApp(appID)
		if err != nil {
			return fiber.NewError(errutil.MapStatusCode(err))
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func enableObservable(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		err := dataStore.SetObservable(appID, true)
		if err != nil {
			return fiber.NewError(errutil.MapStatusCode(err))
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func disableObservable(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		err := dataStore.SetObservable(appID, false)
		if err != nil {
			return fiber.NewError(errutil.MapStatusCode(err))
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
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
