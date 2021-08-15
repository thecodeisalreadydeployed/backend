package apiserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
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
}
