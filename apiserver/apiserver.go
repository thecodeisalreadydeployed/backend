package apiserver

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/workloadcontroller"
	"log"
)

func APIServer(port int) {
	app := fiber.New()

	app.Get("/project/:projectID", func(c *fiber.Ctx) error {
		projectID := cast.ToString(c.Params("projectID"))
		return c.SendString(projectID)
	})

	app.Get("/app/:appID", func(c *fiber.Ctx) error {
		appID := cast.ToString(c.Params("appID"))
		return c.SendString(appID)
	})

	app.Get("/deployment/:deploymentID", func(c *fiber.Ctx) error {
		deploymentID := cast.ToString(c.Params("deploymentID"))
		return c.SendString(deploymentID)
	})

	app.Get("/deployment/:deploymentID/event", func(c *fiber.Ctx) error {
		deploymentID := cast.ToString(c.Params("deploymentID"))
		return c.SendString(deploymentID)
	})

	app.Post("/create", func(c *fiber.Ctx) error {
		payload := model.Payload{}
		if err := c.BodyParser(&payload); err != nil {
			return c.SendStatus(500)
		}
		workloadcontroller.CreateWorkload(&payload)
		return c.SendStatus(200)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
