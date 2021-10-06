package dto

import (
	"github.com/thecodeisalreadydeployed/model"
	"time"
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
}

func (req *CreateAppRequest) ToModel() model.App {
	return model.App{
		ID:        model.GenerateAppID(),
		ProjectID: req.ProjectID,
		Name:      req.Name,
		BuildConfiguration: model.BuildConfiguration{
			BuildScript:      req.BuildScript,
			ParseBuildScript: true,
			InstallCommand:   req.InstallCommand,
			BuildCommand:     req.BuildCommand,
			OutputDirectory:  req.OutputDirectory,
			StartCommand:     req.StartCommand,
		},
		GitSource: model.GitSource{
			RepositoryURL: req.RepositoryURL,
		},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Observable: false,
	}
}
