package apiserver

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/model"
	"log"

	"github.com/thecodeisalreadydeployed/datastore"

	"github.com/gofiber/fiber/v2"
)

func APIServer(port int) {
	app := fiber.New()

	app.Get("/project/get/:projectID", func(c *fiber.Ctx) error {
		projectID := c.Params("projectID")
		result, err := datastore.GetProjectByID(projectID)
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/project/get/:projectID/apps", func(c *fiber.Ctx) error {
		result, err := datastore.GetAppsByProjectID(c.Params("projectID"))
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/app/get/:appID", func(c *fiber.Ctx) error {
		appID := c.Params("appID")
		result, err := datastore.GetAppByID(appID)
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/app/get/:appID/deployments", func(c *fiber.Ctx) error {
		result, err := datastore.GetDeploymentsByAppID(c.Params("appID"))
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/deployment/get/:deploymentID", func(c *fiber.Ctx) error {
		deploymentID := c.Params("deploymentID")
		result, err := datastore.GetDeploymentByID(deploymentID)
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.JSON(result)
	})

	app.Get("/deployment/get/:deploymentID/event", func(c *fiber.Ctx) error {
		event, err := datastore.GetEventByDeploymentID(c.Params("deploymentID"))
		if err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendString(event)
	})

	app.Post("/project/save", func(c *fiber.Ctx) error {
		payload := model.Project{}
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		if err := datastore.SaveProject(&payload); err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/app/save", func(c *fiber.Ctx) error {
		payload := model.App{}
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		if err := datastore.SaveApp(&payload); err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/project/remove/:projectID", func(c *fiber.Ctx) error {
		projectID := c.Params("projectID")
		if err := datastore.RemoveProject(projectID); err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/app/remove/:appID", func(c *fiber.Ctx) error {
		appID := c.Params("appID")
		if err := datastore.RemoveApp(appID); err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
