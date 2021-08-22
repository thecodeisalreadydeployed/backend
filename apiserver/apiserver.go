package apiserver

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/workloadcontroller"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
)

func APIServer(port int) {
	app := fiber.New()

	app.Get("/project/:projectID", func(c *fiber.Ctx) error {
		query := new(model.Project)
		query.ID = c.Params("projectID")
		result := datastore.GetProjectByID(datamodel.NewProjectFromModel(*query))
		return c.JSON(result)
	})

	app.Get("/project/:projectID/apps", func(c *fiber.Ctx) error {
		result := datastore.GetProjectApps(c.Params("projectID"))
		return c.JSON(result)
	})

	app.Get("/app/:appID", func(c *fiber.Ctx) error {
		query := new(model.App)
		query.ID = c.Params("appID")
		result := datastore.GetAppByID(datamodel.NewAppFromModel(*query))
		return c.JSON(result)
	})

	app.Get("/app/:appID/deployments", func(c *fiber.Ctx) error {
		result := datastore.GetAppDeployments(c.Params("appID"))
		return c.JSON(result)
	})

	app.Get("/deployment/:deploymentID", func(c *fiber.Ctx) error {
		query := new(model.Deployment)
		query.ID = c.Params("deploymentID")
		result := datastore.GetDeploymentByID(datamodel.NewDeploymentFromModel(*query))
		return c.JSON(result)
	})

	app.Get("/deployment/:deploymentID/event", func(c *fiber.Ctx) error {
		event := datastore.GetEvent(c.Params("deploymentID"))
		return c.SendString(event)
	})

	app.Post("/project/new", func(c *fiber.Ctx) error {
		payload := dto.CreateProjectRequest{}
		if err := c.BodyParser(&payload); err != nil {
			return c.SendStatus(500)
		}
		return c.SendStatus(200)
	})

	// TODO: Delete this.
	app.Get("/test", func(c *fiber.Ctx) error {
		payload := dto.CreateProjectRequest{
			Name: "test",
		}
		yaml := workloadcontroller.CreateWorkload(&payload)
		return c.SendString(yaml)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
