package datamodel

import (
	"encoding/base64"
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

	buildConfiguration, err := json.Marshal(app.BuildConfiguration)
	if err != nil {
		panic(err)
	}

	buildConfiguration64 := base64.StdEncoding.EncodeToString(buildConfiguration)

	return &App{
		ID:                 app.ID,
		ProjectID:          app.ProjectID,
		Name:               app.Name,
		GitSource:          cast.ToString(gitSource),
		CreatedAt:          app.CreatedAt,
		UpdatedAt:          app.UpdatedAt,
		BuildConfiguration: buildConfiguration64,
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

	buildConfiguration, err := json.Marshal(dpl.BuildConfiguration)
	if err != nil {
		panic(err)
	}

	buildConfiguration64 := base64.StdEncoding.EncodeToString(buildConfiguration)

	return &Deployment{
		ID:                 dpl.ID,
		AppID:              dpl.AppID,
		Meta:               dpl.Meta,
		State:              dpl.State,
		GitSource:          cast.ToString(gitSource),
		Creator:            cast.ToString(creator),
		BuildConfiguration: buildConfiguration64,
		BuiltAt:            dpl.BuiltAt,
		CommittedAt:        dpl.CommittedAt,
		DeployedAt:         dpl.DeployedAt,
		CreatedAt:          dpl.CreatedAt,
		UpdatedAt:          dpl.UpdatedAt,
	}
}
