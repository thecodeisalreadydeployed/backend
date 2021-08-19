package datamodel

import (
	"encoding/json"

	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/model"
)

func NewProjectFromModel(prj model.Project) Project {
	return Project{
		ID:        prj.ID,
		Name:      prj.Name,
		CreatedAt: prj.CreatedAt,
		UpdatedAt: prj.UpdatedAt,
	}
}

func NewAppFromModel(app model.App) App {
	gitSource, err := json.Marshal(app.GitSource)
	if err != nil {
		panic(err)
	}

	return App{
		ID:        app.ID,
		Name:      app.Name,
		GitSource: cast.ToString(gitSource),
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}
}

func NewDeploymentFromModel(dpl model.Deployment) Deployment {
	return Deployment{}
}
