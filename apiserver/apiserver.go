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

	app.Get("/project/:projectID/app/:appID", func(c *fiber.Ctx) error {
		projectID := cast.ToString(c.Params("projectID"))
		appID := cast.ToString(c.Params("appID"))
		return c.SendString(projectID)
	})

	app.Get("/project/:projectID/app/:appID/deployment/:deploymentID", func(c *fiber.Ctx) error {
		projectID := cast.ToString(c.Params("projectID"))
		appID := cast.ToString(c.Params("appID"))
		return c.SendString(projectID)
	})
}
