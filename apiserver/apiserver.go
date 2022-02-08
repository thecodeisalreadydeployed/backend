package apiserver

import (
	"fmt"
	"log"

	"github.com/thecodeisalreadydeployed/apiserver/controller"
	"github.com/thecodeisalreadydeployed/gitapi"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
)

func APIServer(port int, workloadController workloadcontroller.WorkloadController) {
	gitAPIBackend := gitapi.NewGitAPIBackend()
	validator.Init()
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	controller.NewProjectController(app.Group("projects"))
	controller.NewAppController(app.Group("apps"), workloadController)
	controller.NewDeploymentController(app.Group("deployments"), workloadController)

	controller.NewHealthController(app.Group("health"))
	controller.NewPresetController(app.Group("presets"))
	controller.NewBuildScriptController(app.Group("build-script"))
	controller.NewGitAPIController(app.Group("gitapi"), gitAPIBackend)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
