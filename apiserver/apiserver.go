package apiserver

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/apiserver/group"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
)

func APIServer(port int) {
	validator.Init()
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/health", group.Health)

	app.Get("/projects", group.ProjectList)
	app.Get("/project/:projectID", group.ProjectID)
	app.Get("/project/:projectID/apps", group.ProjectApps)
	app.Get("/app/:appID", group.AppID)
	app.Get("/app/:appID/deployments", group.AppDeployments)
	app.Get("/deployment/:deploymentID", group.DeploymentID)
	app.Get("/deployment/:deploymentID/event", group.DeploymentEvent)

	app.Post("/project", group.PostProject)
	app.Post("/app", group.PostApp)
	app.Delete("/project/:projectID", group.DeleteProject)
	app.Delete("/app/:appID", group.DeleteApp)

	app.Get("/preset/:framework", group.Preset)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
