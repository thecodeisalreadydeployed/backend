package apiserver

import (
	"fmt"
	"log"

	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func APIServer(port int) {
	validator.Init()
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		ok := datastore.IsReady()
		return c.JSON(map[string]string{"ok": cast.ToString(ok)})
	})

	app.Get("/projects", func(c *fiber.Ctx) error {
		result, err := datastore.GetAllProjects()
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/project/:projectID", func(c *fiber.Ctx) error {
		projectID := c.Params("projectID")
		result, err := datastore.GetProjectByID(projectID)
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/project/:projectID/apps", func(c *fiber.Ctx) error {
		result, err := datastore.GetAppsByProjectID(c.Params("projectID"))
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/app/:appID", func(c *fiber.Ctx) error {
		appID := c.Params("appID")
		result, err := datastore.GetAppByID(appID)
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/app/:appID/deployments", func(c *fiber.Ctx) error {
		result, err := datastore.GetDeploymentsByAppID(c.Params("appID"))
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/deployment/:deploymentID", func(c *fiber.Ctx) error {
		deploymentID := c.Params("deploymentID")
		result, err := datastore.GetDeploymentByID(deploymentID)
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/deployment/:deploymentID/event", func(c *fiber.Ctx) error {
		event, err := datastore.GetEventByDeploymentID(c.Params("deploymentID"))
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendString(event)
	})

	app.Post("/project", func(c *fiber.Ctx) error {
		request := dto.CreateProjectRequest{}
		if err := c.BodyParser(&request); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}

		if validationErrors := validator.Validate(request); len(validationErrors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
		}

		prj := request.ToModel()

		project, projectErr := datastore.SaveProject(&prj)

		if projectErr != nil {
			return fiber.NewError(mapStatusCode(projectErr))
		}

		return c.JSON(project)
	})

	app.Post("/app", func(c *fiber.Ctx) error {
		request := dto.CreateAppRequest{}

		if err := c.BodyParser(&request); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}

		if validationErrors := validator.Validate(request); len(validationErrors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
		}

		app := request.ToModel()

		if err := datastore.SaveApp(&app); err != nil {
			return fiber.NewError(mapStatusCode(err))
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Delete("/project/:projectID", func(c *fiber.Ctx) error {
		projectID := c.Params("projectID")
		if err := datastore.RemoveProject(projectID); err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendStatus(fiber.StatusOK)
	})

	app.Delete("/app/:appID", func(c *fiber.Ctx) error {
		appID := c.Params("appID")
		if err := datastore.RemoveApp(appID); err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendStatus(fiber.StatusOK)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
