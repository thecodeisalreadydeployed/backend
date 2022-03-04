package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/datastore"
)

func NewPresetController(api fiber.Router, dataStore datastore.DataStore) {
	api.Get("/list", listPresets(dataStore))
	api.Get("/search", searchPreset(dataStore))
	api.Get("/:presetID", getPreset(dataStore))
	api.Post("/", createPreset(dataStore))
	api.Delete("/:presetID", deletePreset(dataStore))
}

func listPresets(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		result, err := dataStore.GetAllPresets()
		return writeResponse(ctx, result, err)
	}
}

func getPreset(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		presetID := ctx.Params("presetID")
		result, err := dataStore.GetPresetByID(presetID)
		return writeResponse(ctx, result, err)
	}
}

func searchPreset(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		presetName := ctx.Query("name")
		result, err := dataStore.GetPresetsByName(presetName)
		return writeResponse(ctx, result, err)
	}
}

func createPreset(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		input := dto.CreatePresetRequest{}

		if err := ctx.BodyParser(&input); err != nil {
			return fiber.NewError(errutil.MapStatusCode(err))
		}

		if validationErrors := validator.CheckStruct(input); len(validationErrors) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErrors)
		}

		inputModel := input.ToModel()
		app, createErr := dataStore.SavePreset(&inputModel)

		return writeResponse(ctx, app, createErr)
	}
}

func deletePreset(dataStore datastore.DataStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		presetID := ctx.Params("presetID")
		err := dataStore.RemovePreset(presetID)
		if err != nil {
			return fiber.NewError(errutil.MapStatusCode(err))
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}
