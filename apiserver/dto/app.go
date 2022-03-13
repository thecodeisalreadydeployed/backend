package dto

import (
	"github.com/thecodeisalreadydeployed/model"
)

type CreateAppRequest struct {
	ProjectID     string `json:"projectID" validate:"required"`
	Name          string `json:"name" validate:"required"`
	RepositoryURL string `json:"repositoryURL" validate:"required"`
	BuildScript   string `json:"buildScript" validate:"required"`
	Branch        string `json:"branch" validate:"required"`
	FetchInterval int    `json:"fetchInterval"`
}

func (req *CreateAppRequest) ToModel() model.App {
	return model.App{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		BuildConfiguration: model.BuildConfiguration{
			BuildScript: req.BuildScript,
		},
		GitSource: model.GitSource{
			RepositoryURL: req.RepositoryURL,
			Branch:        req.Branch,
		},
		Observable:    true,
		FetchInterval: req.FetchInterval,
	}
}
