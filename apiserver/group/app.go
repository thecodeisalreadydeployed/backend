package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
)

func AppID(c *fiber.Ctx) error {
	appID := c.Params("appID")
	result, err := datastore.GetAppByID(appID)
	if err != nil {
		return fiber.NewError(apiserver.MapStatusCode(err))
	}
	return c.JSON(result)
}

func AppDeployments(c *fiber.Ctx) error {
	result, err := datastore.GetDeploymentsByAppID(c.Params("appID"))
	if err != nil {
		return fiber.NewError(apiserver.MapStatusCode(err))
	}
	return c.JSON(result)
}

func PostApp(c *fiber.Ctx) error {
	request := dto.CreateAppRequest{}

	if err := c.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	if validationErrors := validator.CheckStruct(request); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	app := request.ToModel()

	if err := datastore.SaveApp(&app); err != nil {
		return fiber.NewError(apiserver.MapStatusCode(err))
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteApp(c *fiber.Ctx) error {
	appID := c.Params("appID")
	if err := datastore.RemoveApp(appID); err != nil {
		return fiber.NewError(apiserver.MapStatusCode(err))
	}
	return c.SendStatus(fiber.StatusOK)
}
