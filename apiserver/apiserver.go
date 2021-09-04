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

	app.Post("/project/new", func(c *fiber.Ctx) error {
		payload := model.Project{}
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		if err := datastore.CreateProject(&payload); err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendStatus(200)
	})

	app.Post("/app/new", func(c *fiber.Ctx) error {
		payload := model.App{}
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		if err := datastore.CreateApp(&payload); err != nil {
			return fiber.NewError(mapStatusCode(err))
		}
		return c.SendStatus(200)
	})

	app.Post("/project/update", func(c *fiber.Ctx) error {
		payload := model.Project{}
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		return c.SendStatus(200)
	})

	app.Post("/app/update", func(c *fiber.Ctx) error {
		payload := model.App{}
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest)
		}
		return c.SendStatus(200)
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusOK)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
