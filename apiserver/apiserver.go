package apiserver

import (
	"fmt"
	"log"

	"github.com/thecodeisalreadydeployed/repositoryobserver"
	"github.com/thecodeisalreadydeployed/statusapi"

	"github.com/thecodeisalreadydeployed/apiserver/auth"
	"github.com/thecodeisalreadydeployed/apiserver/controller"
	"github.com/thecodeisalreadydeployed/gitapi"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
)

func APIServer(port int, workloadController workloadcontroller.WorkloadController, observer repositoryobserver.RepositoryObserver) {
	firebaseAuth := auth.SetupFirebase()
	gitAPIBackend := gitapi.NewGitAPIBackend(zap.L())
	statusAPIBackend := statusapi.NewStatusAPIBackend(zap.L())
	validator.Init()
	app := fiber.New()

	app.Use(cors.New())

	controller.NewProjectController(app.Group("projects", auth.EnsureAuthenticated(firebaseAuth)), workloadController)
	controller.NewAppController(app.Group("apps", auth.EnsureAuthenticated(firebaseAuth)), workloadController, observer, statusAPIBackend)
	controller.NewDeploymentController(app.Group("deployments", auth.EnsureAuthenticated(firebaseAuth)), workloadController)
	controller.NewPresetController(app.Group("presets", auth.EnsureAuthenticated(firebaseAuth)))
	controller.NewGitAPIController(app.Group("gitapi", auth.EnsureAuthenticated(firebaseAuth)), gitAPIBackend)

	controller.NewHealthController(app.Group("health"))

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
