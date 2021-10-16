package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
)

func ProjectList(c *fiber.Ctx) error {
	result, err := datastore.GetAllProjects()
	if err != nil {
		return fiber.NewError(apiserver.MapStatusCode(err))
	}
	return c.JSON(result)
}

func ProjectID(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	result, err := datastore.GetProjectByID(projectID)
	if err != nil {
		return fiber.NewError(apiserver.MapStatusCode(err))
	}
	return c.JSON(result)
}

func ProjectApps(c *fiber.Ctx) error {
	result, err := datastore.GetAppsByProjectID(c.Params("projectID"))
	if err != nil {
		return fiber.NewError(apiserver.MapStatusCode(err))
	}
	return c.JSON(result)
}

func PostProject(c *fiber.Ctx) error {
	request := dto.CreateProjectRequest{}
	if err := c.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	if validationErrors := validator.CheckStruct(request); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	prj := request.ToModel()

	project, projectErr := datastore.SaveProject(&prj)

	if projectErr != nil {
		return fiber.NewError(apiserver.MapStatusCode(projectErr))
	}

	return c.JSON(project)
}

func DeleteProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	if err := datastore.RemoveProject(projectID); err != nil {
		return fiber.NewError(apiserver.MapStatusCode(err))
	}
	return c.SendStatus(fiber.StatusOK)
}
