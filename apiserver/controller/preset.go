package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewPresetController(api fiber.Router) {
	api.Get("/list", listPresets)
	api.Get("/search", searchPreset)
	api.Get("/:presetID", getPreset)
	api.Post("/", createPreset)
	api.Delete("/:presetID", deletePreset)
}

func listPresets(ctx *fiber.Ctx) error {
	result, err := datastore.GetAllPresets(datastore.GetDB())
	return writeResponse(ctx, result, err)
}

func getPreset(ctx *fiber.Ctx) error {
	presetID := ctx.Params("presetID")
	result, err := datastore.GetPresetByID(datastore.GetDB(), presetID)
	return writeResponse(ctx, result, err)
}

func searchPreset(ctx *fiber.Ctx) error {
	presetName := ctx.Query("name")
	result, err := datastore.GetPresetsByName(datastore.GetDB(), presetName)
	return writeResponse(ctx, result, err)
}

func createPreset(ctx *fiber.Ctx) error {
	input := dto.CreatePresetRequest{}

	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}

	if validationErrors := validator.CheckStruct(input); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	inputModel := input.ToModel()
	app, createErr := datastore.SavePreset(datastore.GetDB(), &inputModel)

	return writeResponse(ctx, app, createErr)
}

func deletePreset(ctx *fiber.Ctx) error {
	presetID := ctx.Params("presetID")
	err := datastore.RemovePreset(datastore.GetDB(), presetID)
	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.SendStatus(fiber.StatusOK)
}
