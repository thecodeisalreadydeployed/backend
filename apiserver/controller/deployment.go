package controller

import (
	"github.com/thecodeisalreadydeployed/repositoryobserver"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"
	"go.uber.org/zap"
)

func NewDeploymentController(api fiber.Router, workloadController workloadcontroller.WorkloadController, observer repositoryobserver.RepositoryObserver) {
	// Create a new deployment
	api.Post("/", createDeployment(workloadController))

	api.Get("/:deploymentID", getDeployment)
	api.Get("/:deploymentID/events", getDeploymentEvents)
	api.Post("/:deploymentID/events", createDeploymentEvents)
	api.Delete("/:deploymentID", deleteDeployment)
	api.Get("/:deploymentID/refresh", forceRefresh(observer))
}

func createDeployment(workloadController workloadcontroller.WorkloadController) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		deployment, err := workloadController.NewDeployment(appID, nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}
		return ctx.JSON(deployment)
	}
}

func getDeployment(ctx *fiber.Ctx) error {
	deploymentID := ctx.Params("deploymentID")
	result, err := datastore.GetDeploymentByID(datastore.GetDB(), deploymentID)
	return writeResponse(ctx, result, err)
}

func getDeploymentEvents(ctx *fiber.Ctx) error {
	deploymentID := ctx.Params("deploymentID")
	result, err := datastore.GetEventsByDeploymentID(datastore.GetDB(), deploymentID)
	if err != nil {
		return writeResponse(ctx, result, err)
	}

	ret := *result
	sort.SliceStable(ret, func(i, j int) bool {
		// a, _ := ksuid.Parse(ret[i].ID)
		// b, _ := ksuid.Parse(ret[j].ID)
		// return ksuid.Compare(a, b) < 0
		return ret[i].ExportedAt.Before(ret[j].ExportedAt)
	})

	return writeResponse(ctx, ret, nil)
}

func createDeploymentEvents(ctx *fiber.Ctx) error {
	isInternalRequest := string(ctx.Request().Header.Peek("X-CodeDeploy-Internal-Request")) == "True" && len(ctx.Request().Header.Peek("X-Forwarded-For")) == 0
	if !isInternalRequest {
		return fiber.NewError(fiber.StatusNotFound)
	}

	deploymentID := ctx.Params("deploymentID")
	input := dto.CreateDeploymentEventRequest{}

	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}

	if validationErrors := validator.CheckStruct(input); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	deployment, err := datastore.GetDeploymentByID(datastore.GetDB(), deploymentID)
	if err != nil {
		return err
	}

	events, err := datastore.GetEventsByDeploymentID(datastore.GetDB(), deployment.ID)
	if err != nil {
		return err
	}

	if len(*events) == 0 && deployment.State == model.DeploymentStateQueueing {
		err = datastore.SetDeploymentState(datastore.GetDB(), deployment.ID, model.DeploymentStateBuilding)
		if err != nil {
			return err
		}
	}

	inputModel := input.ToModel()
	inputModel.DeploymentID = deployment.ID

	zap.L().Sugar().Debug(inputModel.Text)

	event, createErr := datastore.SaveEvent(datastore.GetDB(), &inputModel)
	return writeResponse(ctx, event, createErr)
}

func deleteDeployment(ctx *fiber.Ctx) error {
	deploymentID := ctx.Params("deploymentID")
	err := datastore.RemoveDeployment(datastore.GetDB(), deploymentID)
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func forceRefresh(observer repositoryobserver.RepositoryObserver) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		zap.L().Warn(`Depending on the circumstances, refreshing may not give you the desired results.
Refreshing to fetch the codebase will only skip the idle duration of the repository observer.
If the observer is refreshed while it is deploying, the observer will immediately fetch the codebase after deployment.
If the observer is refreshed while waiting after an error occurs, the observer will not restart the deployment.
If the observer is refreshed twice while idle, it will deploy the codebase twice.
Please use the refresh function with full understanding and only when the observer is idle.`)
		deploymentID := ctx.Params("deploymentID")
		dpl, err := datastore.GetDeploymentByID(datastore.GetDB(), deploymentID)
		if err != nil {
			return fiber.NewError(errutil.MapStatusCode(err))
		}
		observer.Refresh(dpl.AppID)
		return ctx.SendStatus(fiber.StatusOK)
	}
}
