package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/apiserver/errutil"
	"github.com/thecodeisalreadydeployed/apiserver/validator"
	"github.com/thecodeisalreadydeployed/gitapi"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"strings"
)

func NewGitApiController(api fiber.Router) {
	api.Post("/branches", getBranches)
	api.Post("/files", getFiles)
	api.Post("/raw", getRaw)
}

func getBranches(ctx *fiber.Ctx) error {
	input := dto.GetBranchesRequest{}

	if err := validator.ParseBodyAndValidate(ctx, &input); err != nil {
		return err
	}

	var branches []string
	var err error
	if isGitHubURL(input.Url) {
		branches, err = gitapi.GetBranches(input.Url)
	} else {
		git, gitErr := gitgateway.NewGitGatewayRemote(input.Url)
		if gitErr != nil {
			return fiber.ErrBadRequest
		}
		branches, err = git.GetBranches()
	}

	return writeResponse(ctx, branches, err)
}

func getFiles(ctx *fiber.Ctx) error {
	input := dto.GetFilesRequest{}

	if err := validator.ParseBodyAndValidate(ctx, &input); err != nil {
		return err
	}

	var files []string
	var err error
	if isGitHubURL(input.Url) {
		files, err = gitapi.GetFiles(input.Url, input.Branch)
	} else {
		git, gitErr := gitgateway.NewGitGatewayRemote(input.Url)
		if gitErr != nil {
			return fiber.ErrBadRequest
		}
		files, err = git.GetFiles(input.Branch)
	}

	return writeResponse(ctx, files, err)
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
