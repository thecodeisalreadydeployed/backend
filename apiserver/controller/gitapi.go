package controller

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/gitapi"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
)

func NewGitAPIController(api fiber.Router, gitAPIBackend gitapi.GitAPIBackend) {
	api.Post("/branches", getBranches(gitAPIBackend))
	api.Post("/files", getFiles(gitAPIBackend))
	api.Post("/raw", getRaw)
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

func getRaw(ctx *fiber.Ctx) error {
	input := dto.GetRawRequest{}

	if err := validator.ParseBodyAndValidate(ctx, &input); err != nil {
		return err
	}

	var raw string
	var err error
	if isGitHubURL(input.Url) {
		raw, err = gitapi.GetRaw(input.Url, input.Branch, input.Path)
	} else {
		git, gitErr := gitgateway.NewGitGatewayRemote(input.Url)
		if gitErr != nil {
			return fiber.ErrBadRequest
		}
		raw, err = git.GetRaw(input.Branch, input.Path)
	}

	if err != nil {
		return fiber.NewError(errutil.MapStatusCode(err))
	}
	return ctx.SendString(raw)
}

func isGitHubURL(url string) bool {
	slices := strings.Split(url, "/")
	return len(slices) == 5 && slices[2] == "github.com"
}
