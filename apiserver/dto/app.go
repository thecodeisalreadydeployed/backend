package dto

import (
	"github.com/thecodeisalreadydeployed/model"
)

type CreateAppRequest struct {
	ProjectID       string `validate:"required"`
	Name            string `validate:"required"`
	RepositoryURL   string `validate:"required"`
	BuildScript     string `validate:"required"`
	InstallCommand  string `validate:"required"`
	BuildCommand    string `validate:"required"`
	OutputDirectory string `validate:"required"`
	StartCommand    string `validate:"required"`
	Branch          string `validate:"required"`
}

func (req *CreateAppRequest) ToModel() model.App {
	return model.App{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		BuildConfiguration: model.BuildConfiguration{
			BuildScript:     req.BuildScript,
			InstallCommand:  req.InstallCommand,
			BuildCommand:    req.BuildCommand,
			OutputDirectory: req.OutputDirectory,
			StartCommand:    req.StartCommand,
		},
		GitSource: model.GitSource{
			RepositoryURL: req.RepositoryURL,
			Branch:        req.Branch,
		},
		Observable: false,
	}
}
