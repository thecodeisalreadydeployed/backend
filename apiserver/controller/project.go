package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewProjectController(api fiber.Router) {
	api.Get("/", listProjects)
	api.Get("/:projectID", getProject)
	api.Get("/:projectID/apps", listProjectApps)
}

func listProjects(ctx *fiber.Ctx) error {
	result, err := datastore.GetAllProjects()
	return writeResponse(ctx, result, err)
}

func getProject(ctx *fiber.Ctx) error {
	projectID := ctx.Params("projectID")
	result, err := datastore.GetProjectByID(projectID)
	return writeResponse(ctx, result, err)
}

func listProjectApps(ctx *fiber.Ctx) error {
	projectID := ctx.Params("projectID")
	result, err := datastore.GetAppsByProjectID(projectID)
	return writeResponse(ctx, result, err)
}
