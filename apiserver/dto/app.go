package dto

import "github.com/thecodeisalreadydeployed/model"

type CreateAppRequest struct {
	ProjectID       string `validate:"required"`
	Name            string `validate:"required"`
	RepositoryURL   string `validate:"required"`
	BuildCommand    string `validate:"required"`
	OutputDirectory string `validate:"required"`
	InstallCommand  string `validate:"required"`
	StartCommand    string `validate:"required"`
}

func (req *CreateAppRequest) ToModel() model.App {
	return model.App{
		ProjectID:       req.ProjectID,
		Name:            req.Name,
		BuildCommand:    req.BuildCommand,
		OutputDirectory: req.OutputDirectory,
		InstallCommand:  req.InstallCommand,
		GitSource: model.GitSource{
			RepositoryURL: req.RepositoryURL,
		},
	}
}
