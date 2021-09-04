package datamodel

import (
	"encoding/json"

	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/model"
)

func NewProjectFromModel(prj *model.Project) *Project {
	return &Project{
		ID:        prj.ID,
		Name:      prj.Name,
		CreatedAt: prj.CreatedAt,
		UpdatedAt: prj.UpdatedAt,
	}
}

func NewAppFromModel(app *model.App) *App {
	gitSource, err := json.Marshal(app.GitSource)
	if err != nil {
		panic(err)
	}

	return &App{
		ID:              app.ID,
		ProjectID:       app.ProjectID,
		Name:            app.Name,
		GitSource:       cast.ToString(gitSource),
		CreatedAt:       app.CreatedAt,
		UpdatedAt:       app.UpdatedAt,
		BuildCommand:    app.BuildCommand,
		OutputDirectory: app.OutputDirectory,
		InstallCommand:  app.InstallCommand,
	}
}

func NewDeploymentFromModel(dpl *model.Deployment) *Deployment {
	gitSource, err := json.Marshal(dpl.GitSource)
	if err != nil {
		panic(err)
	}

	creator, err := json.Marshal(dpl.Creator)
	if err != nil {
		panic(err)
	}

	return &Deployment{
		ID:        dpl.ID,
		AppID:     dpl.AppID,
		Name:      dpl.Name,
		Meta:      dpl.Meta,
		State:     dpl.State,
		GitSource: cast.ToString(gitSource),
		Creator:   cast.ToString(creator),
		BuiltAt:   dpl.BuiltAt,
		CreatedAt: dpl.CreatedAt,
		UpdatedAt: dpl.UpdatedAt,
	}
}
