package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewProjectController(api fiber.Router) {
	api.Get("/list", listProjects)
	api.Get("/:projectID", getProject)
	api.Get("/:projectID/apps", listProjectApps)
	api.Get("/name/:projectName", searchProject)
	api.Post("/", createProject)
	api.Delete("/:projectID", deleteProject)
}

func listProjects(ctx *fiber.Ctx) error {
	result, err := datastore.GetAllProjects(datastore.GetDB())
	return writeResponse(ctx, result, err)
}

func getProject(ctx *fiber.Ctx) error {
	projectID := ctx.Params("projectID")
	result, err := datastore.GetProjectByID(datastore.GetDB(), projectID)
	return writeResponse(ctx, result, err)
}

func listProjectApps(ctx *fiber.Ctx) error {
	projectID := ctx.Params("projectID")
	result, err := datastore.GetAppsByProjectID(datastore.GetDB(), projectID)
	return writeResponse(ctx, result, err)
}

func searchProject(ctx *fiber.Ctx) error {
	projectName := ctx.Params("projectName")
	result, err := datastore.GetProjectsByName(datastore.GetDB(), projectName)
	return writeResponse(ctx, result, err)
}

func createProject(ctx *fiber.Ctx) error {
	input := dto.CreateProjectRequest{}

	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}

	if validationErrors := validator.CheckStruct(input); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	prj := input.ToModel()
	project, createErr := datastore.SaveProject(datastore.GetDB(), &prj)

	return writeResponse(ctx, project, createErr)
}

func deleteProject(ctx *fiber.Ctx) error {
	projectID := ctx.Params("projectID")
	err := datastore.RemoveProject(datastore.GetDB(), projectID)
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.SendStatus(fiber.StatusOK)
}
