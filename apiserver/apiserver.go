package apiserver

import (
	"fmt"
	"log"

	"github.com/thecodeisalreadydeployed/apiserver/controller"

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

	controller.NewProjectController(app.Group("projects"))

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
