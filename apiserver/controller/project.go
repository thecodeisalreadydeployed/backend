package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"
)

func NewProjectController(api fiber.Router, workloadController workloadcontroller.WorkloadController) {
	api.Get("/list", listProjects)
	api.Get("/search", searchProject)
	api.Get("/:projectID", getProject)
	api.Get("/:projectID/apps", listProjectApps)
	api.Post("/", createProject(workloadController))
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
	projectName := ctx.Query("name")
	result, err := datastore.GetProjectsByName(datastore.GetDB(), projectName)
	return writeResponse(ctx, result, err)
}

func createProject(workloadController workloadcontroller.WorkloadController) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := dto.CreateProjectRequest{}
		if err := validator.ParseBodyAndValidate(c, &input); err != nil {
			return err
		}
		prj := input.ToModel()
		project, createErr := workloadController.NewProject(&prj)
		return writeResponse(c, project, createErr)
	}

}

func deleteProject(ctx *fiber.Ctx) error {
	projectID := ctx.Params("projectID")
	err := datastore.RemoveProject(datastore.GetDB(), projectID)
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.SendStatus(fiber.StatusOK)
}
