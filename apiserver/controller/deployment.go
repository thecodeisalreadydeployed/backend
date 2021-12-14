package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/datastore"
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
	return fiber.NewError(fiber.StatusNotImplemented)
}

func createDeploymentEvents(ctx *fiber.Ctx) error {
	isInternalRequest := string(ctx.Request().Header.Peek("X-CodeDeploy-Internal-Request")) == "True" && len(ctx.Request().Header.Peek("X-Forwarded-For")) == 0
	if !isInternalRequest {
		return fiber.NewError(fiber.StatusNotFound)
	}

	_ = ctx.Params("deploymentID")
	return fiber.NewError(fiber.StatusNotImplemented)
}
