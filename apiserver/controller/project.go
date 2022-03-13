package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"
)

func NewProjectController(api fiber.Router, workloadController workloadcontroller.WorkloadController, dataStore datastore.DataStore) {
	api.Get("/list", listProjects(dataStore))
	api.Get("/search", searchProject(dataStore))
	api.Get("/:projectID", getProject(dataStore))
	api.Get("/:projectID/apps", listProjectApps(dataStore))
	api.Post("/", createProject(workloadController, dataStore))
	api.Delete("/:projectID", deleteProject(dataStore))
}

func listProjects(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		result, err := dataStore.GetAllProjects()
		return writeResponse(ctx, result, err)
	}
}

func getProject(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		projectID := ctx.Params("projectID")
		result, err := dataStore.GetProjectByID(projectID)
		return writeResponse(ctx, result, err)
	}
}

func listProjectApps(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		projectID := ctx.Params("projectID")
		result, err := dataStore.GetAppsByProjectID(projectID)
		return writeResponse(ctx, result, err)
	}
}

func searchProject(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		projectName := ctx.Query("name")
		result, err := dataStore.GetProjectsByName(projectName)
		return writeResponse(ctx, result, err)
	}
}

func createProject(workloadController workloadcontroller.WorkloadController, dataStore datastore.DataStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := dto.CreateProjectRequest{}
		if err := validator.ParseBodyAndValidate(c, &input); err != nil {
			return err
		}
		prj := input.ToModel()
		project, createErr := workloadController.NewProject(&prj, dataStore)
		return writeResponse(c, project, createErr)
	}

}

func deleteProject(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		projectID := ctx.Params("projectID")
		err := dataStore.RemoveProject(projectID)
		if err != nil {
			return fiber.NewError(errutil.MapStatusCode(err))
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}
