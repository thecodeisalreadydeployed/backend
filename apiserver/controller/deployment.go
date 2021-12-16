package controller

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
	"go.uber.org/zap"
)

func NewDeploymentController(api fiber.Router) {
	api.Get("/:deploymentID", getDeployment)
	api.Get("/:deploymentID/events", getDeploymentEvents)
	api.Post("/:deploymentID/events", createDeploymentEvents)
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
		a, _ := ksuid.Parse(ret[i].ID)
		b, _ := ksuid.Parse(ret[j].ID)
		return ksuid.Compare(a, b) < 0
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

	inputModel := input.ToModel()
	inputModel.DeploymentID = deploymentID

	zap.L().Sugar().Debug(inputModel.Text)

	_ = datastore.SaveEvent(datastore.GetDB(), &inputModel)

	return ctx.SendStatus(fiber.StatusOK)
}
