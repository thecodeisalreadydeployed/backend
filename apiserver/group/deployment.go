package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/datastore"
)

func DeploymentID(c *fiber.Ctx) error {
	deploymentID := c.Params("deploymentID")
	result, err := datastore.GetDeploymentByID(deploymentID)
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return c.JSON(result)
}

func DeploymentEvent(c *fiber.Ctx) error {
	event, err := datastore.GetEventByDeploymentID(c.Params("deploymentID"))
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return c.SendString(event)
}
