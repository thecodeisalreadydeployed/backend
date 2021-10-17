package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewProjectController(api fiber.Router) {
	api.Get("/", listProjects)
}

func listProjects(ctx *fiber.Ctx) error {
	result, err := datastore.GetAllProjects()
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.JSON(result)
}

func getProject(ctx *fiber.Ctx) error {
	projectID := ctx.Params("projectID")
	result, err := datastore.GetProjectByID(projectID)
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.JSON(result)
}
