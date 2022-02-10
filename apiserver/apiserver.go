package apiserver

import (
	"fmt"
	"log"

	"github.com/thecodeisalreadydeployed/apiserver/auth"
	"github.com/thecodeisalreadydeployed/apiserver/controller"
	"github.com/thecodeisalreadydeployed/gitapi"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
)

func APIServer(port int, workloadController workloadcontroller.WorkloadController) {
	gitAPIBackend := gitapi.NewGitAPIBackend(zap.L())
	validator.Init()
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	controller.NewProjectController(app.Group("projects", auth.EnsureAuthenticated()), workloadController)
	controller.NewAppController(app.Group("apps", auth.EnsureAuthenticated()), workloadController)
	controller.NewDeploymentController(app.Group("deployments", auth.EnsureAuthenticated()), workloadController)
	controller.NewPresetController(app.Group("presets", auth.EnsureAuthenticated()))
	controller.NewGitAPIController(app.Group("gitapi", auth.EnsureAuthenticated()), gitAPIBackend)

	controller.NewHealthController(app.Group("health"))

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
