package apiserver

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
)

func APIServer(port int) {
	app := fiber.New()

	app.Get("/project/:projectID", func(c *fiber.Ctx) error {
		query := new(model.Project)
		query.ID = c.Params("projectID")
		result := datastore.GetProject(query)
		return c.JSON(result)
	})

	app.Get("/app/:appID", func(c *fiber.Ctx) error {
		query := new(model.App)
		query.ID = c.Params("appID")
		result := datastore.GetApp(query)
		return c.JSON(result)
	})

	app.Get("/deployment/:deploymentID", func(c *fiber.Ctx) error {
		query := new(model.Deployment)
		query.ID = c.Params("appID")
		result := datastore.GetDeployment(query)
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

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
