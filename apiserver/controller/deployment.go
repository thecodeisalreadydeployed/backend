package controller

import (
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

func NewDeploymentController(api fiber.Router, workloadController workloadcontroller.WorkloadController, dataStore datastore.DataStore) {
	// Create a new deployment
	api.Post("/", createDeployment(workloadController, dataStore))

	api.Get("/:deploymentID", getDeployment(dataStore))
	api.Get("/:deploymentID/events", getDeploymentEvents(dataStore))
	api.Post("/:deploymentID/events", createDeploymentEvents(dataStore))
	api.Delete("/:deploymentID", deleteDeployment(dataStore))
}

func createDeployment(workloadController workloadcontroller.WorkloadController, dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		appID := ctx.Params("appID")
		deployment, err := workloadController.NewDeployment(appID, nil, dataStore)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}
		return ctx.JSON(deployment)
	}
}

func getDeployment(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		deploymentID := ctx.Params("deploymentID")
		result, err := dataStore.GetDeploymentByID(deploymentID)
		return writeResponse(ctx, result, err)
	}
}

func getDeploymentEvents(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		deploymentID := ctx.Params("deploymentID")
		result, err := dataStore.GetEventsByDeploymentID(deploymentID)
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
}

func createDeploymentEvents(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
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

		deployment, err := dataStore.GetDeploymentByID(deploymentID)
		if err != nil {
			return err
		}

		events, err := dataStore.GetEventsByDeploymentID(deployment.ID)
		if err != nil {
			return err
		}

		if len(*events) == 0 && deployment.State == model.DeploymentStateQueueing {
			err = dataStore.SetDeploymentState(deployment.ID, model.DeploymentStateBuilding)
			if err != nil {
				return err
			}
		}

		inputModel := input.ToModel()
		inputModel.DeploymentID = deployment.ID

		zap.L().Sugar().Debug(inputModel.Text)

		event, createErr := dataStore.SaveEvent(&inputModel)
		return writeResponse(ctx, event, createErr)
	}
}

func deleteDeployment(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		deploymentID := ctx.Params("deploymentID")
		err := dataStore.RemoveDeployment(deploymentID)
		if err != nil {
			return fiber.NewError(errutil.MapStatusCode(err))
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}
