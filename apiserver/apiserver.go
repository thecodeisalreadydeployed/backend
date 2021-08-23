package apiserver

import (
	"fmt"
	"log"

	"github.com/thecodeisalreadydeployed/workloadcontroller"

	"github.com/thecodeisalreadydeployed/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
)

func APIServer(port int) {
	app := fiber.New()

	app.Get("/project/:projectID", func(c *fiber.Ctx) error {
		projectID := c.Params("projectID")
		result, err := datastore.GetProjectByID(projectID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound)
		}
		return c.JSON(result)
	})

	app.Get("/project/:projectID/apps", func(c *fiber.Ctx) error {
		result := datastore.GetAppsByProjectID(c.Params("projectID"))
		return c.JSON(result)
	})

	app.Get("/app/:appID", func(c *fiber.Ctx) error {
		appID := c.Params("appID")
		result := datastore.GetAppByID(appID)
		return c.JSON(result)
	})

	app.Get("/app/:appID/deployments", func(c *fiber.Ctx) error {
		result := datastore.GetDeploymentsByAppID(c.Params("appID"))
		return c.JSON(result)
	})

	app.Get("/deployment/:deploymentID", func(c *fiber.Ctx) error {
		deploymentID := c.Params("deploymentID")
		result := datastore.GetDeploymentByID(deploymentID)
		return c.JSON(result)
	})

	app.Get("/deployment/:deploymentID/event", func(c *fiber.Ctx) error {
		event := datastore.GetEventByDeploymentID(c.Params("deploymentID"))
		return c.SendString(event)
	})

	app.Post("/project/new", func(c *fiber.Ctx) error {
		payload := dto.CreateProjectRequest{}
		if err := c.BodyParser(&payload); err != nil {
			return c.SendStatus(500)
		}
		return c.SendStatus(200)
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// TODO: Jean will delete this. Don't delete it now.
	app.Get("/test", func(c *fiber.Ctx) error {
		payload := dto.CreateProjectRequest{
			Name: "test",
		}
		yaml := workloadcontroller.CreateWorkload(&payload)
		return c.SendString(yaml)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
