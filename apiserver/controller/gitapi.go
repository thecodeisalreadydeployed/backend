package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/gitapi"
)

func NewGitAPIController(api fiber.Router, gitAPIBackend gitapi.GitAPIBackend) {
	api.Post("/branches", getBranches(gitAPIBackend))
	api.Post("/files", getFiles(gitAPIBackend))
	api.Post("/raw", getRaw(gitAPIBackend))
}

func getBranches(gitAPIBackend gitapi.GitAPIBackend) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := dto.GetBranchesRequest{}

		if err := validator.ParseBodyAndValidate(c, &input); err != nil {
			return err
		}

		branches, err := gitAPIBackend.GetBranches(input.URL)
		return writeResponse(c, branches, err)
	}
}

func getFiles(gitAPIBackend gitapi.GitAPIBackend) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := dto.GetFilesRequest{}

		if err := validator.ParseBodyAndValidate(c, &input); err != nil {
			return err
		}
		files, err := gitAPIBackend.GetFiles(input.URL, input.Branch)
		return writeResponse(c, files, err)
	}
}

func getRaw(gitAPIBackend gitapi.GitAPIBackend) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := dto.GetRawRequest{}

		if err := validator.ParseBodyAndValidate(c, &input); err != nil {
			return err
		}

		raw, err := gitAPIBackend.GetRaw(input.URL, input.Branch, input.Path)
		if err != nil {
			return fiber.NewError(errutil.MapStatusCode(err))
		}
		return c.SendString(raw)
	}
}
