package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewDeploymentController(api fiber.Router) {
	api.Get("/:deploymentID", getDeployment)
	api.Get("/:deploymentID/event", getDeploymentEvent)
}

func getDeployment(ctx *fiber.Ctx) error {
	deploymentID := ctx.Params("deploymentID")
	result, err := datastore.GetDeploymentByID(deploymentID)
	return writeResponse(ctx, result, err)
}

func getDeploymentEvent(ctx *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotImplemented)
}
