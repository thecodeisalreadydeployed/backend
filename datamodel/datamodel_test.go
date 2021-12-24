package datamodel

import (
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/model"
	"testing"
	"time"
)

func TestProject_ToModel(t *testing.T) {
	prj := Project{
		ID:        "prj_test",
		Name:      "Best Project",
		CreatedAt: time.Unix(1, 0),
		UpdatedAt: time.Unix(2, 0),
	}

	expected := model.Project{
		ID:        "prj_test",
		Name:      "Best Project",
		CreatedAt: time.Unix(1, 0),
		UpdatedAt: time.Unix(2, 0),
	}

	actual := prj.ToModel()
	assert.Equal(t, expected, actual)
}

func TestNewProjectFromModel(t *testing.T) {
	prj := model.Project{
		ID:        "prj_test",
		Name:      "Best Project",
		CreatedAt: time.Unix(1, 0),
		UpdatedAt: time.Unix(2, 0),
	}

	expected := &Project{
		ID:        "prj_test",
		Name:      "Best Project",
		CreatedAt: time.Unix(1, 0),
		UpdatedAt: time.Unix(2, 0),
	}

	actual := NewProjectFromModel(&prj)
	assert.Equal(t, expected, actual)
}

func TestApp_ToModel(t *testing.T) {
	app := App{
		ID:                 "app_test",
		ProjectID:          "prj_test",
		Project:            Project{},
		Name:               "Best App",
		GitSource:          model.GetGitSourceString(model.GitSource{}),
		CreatedAt:          time.Unix(1, 0),
		UpdatedAt:          time.Unix(2, 0),
		BuildConfiguration: model.GetBuildConfigurationString(model.BuildConfiguration{}),
		Observable:         false,
	}

	expected := model.App{
		ID:                 "app_test",
		ProjectID:          "prj_test",
		Name:               "Best App",
		GitSource:          model.GitSource{},
		CreatedAt:          time.Unix(1, 0),
		UpdatedAt:          time.Unix(2, 0),
		BuildConfiguration: model.BuildConfiguration{},
		Observable:         false,
	}

	actual := app.ToModel()
	assert.Equal(t, expected, actual)
}

func TestNewAppFromModel(t *testing.T) {
	app := model.App{
		ID:                 "app_test",
		ProjectID:          "prj_test",
		Name:               "Best App",
		GitSource:          model.GitSource{},
		CreatedAt:          time.Unix(1, 0),
		UpdatedAt:          time.Unix(2, 0),
		BuildConfiguration: model.BuildConfiguration{},
		Observable:         false,
	}

	expected := &App{
		ID:                 "app_test",
		ProjectID:          "prj_test",
		Project:            Project{},
		Name:               "Best App",
		GitSource:          model.GetGitSourceString(model.GitSource{}),
		CreatedAt:          time.Unix(1, 0),
		UpdatedAt:          time.Unix(2, 0),
		BuildConfiguration: model.GetBuildConfigurationString(model.BuildConfiguration{}),
		Observable:         false,
	}

	actual := NewAppFromModel(&app)
	assert.Equal(t, expected, actual)
}

func TestDeployment_ToModel(t *testing.T) {
	dpl := Deployment{
		ID:                 "dpl_test",
		AppID:              "app_test",
		App:                App{},
		Creator:            model.GetCreatorString(model.Actor{}),
		Meta:               "dummy_meta",
		GitSource:          model.GetGitSourceString(model.GitSource{}),
		BuiltAt:            time.Unix(0, 0),
		CommittedAt:        time.Unix(0, 0),
		DeployedAt:         time.Unix(0, 0),
		BuildConfiguration: model.GetBuildConfigurationString(model.BuildConfiguration{}),
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		State:              model.DeploymentStateReady,
	}

	expected := model.Deployment{
		ID:                 "dpl_test",
		AppID:              "app_test",
		Creator:            model.Actor{},
		Meta:               "dummy_meta",
		GitSource:          model.GitSource{},
		BuiltAt:            time.Unix(0, 0),
		CommittedAt:        time.Unix(0, 0),
		DeployedAt:         time.Unix(0, 0),
		BuildConfiguration: model.BuildConfiguration{},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		State:              model.DeploymentStateReady,
	}

	actual := dpl.ToModel()
	assert.Equal(t, expected, actual)
}

func TestNewDeploymentFromModel(t *testing.T) {
	dpl := model.Deployment{
		ID:                 "dpl_test",
		AppID:              "app_test",
		Creator:            model.Actor{},
		Meta:               "dummy_meta",
		GitSource:          model.GitSource{},
		BuiltAt:            time.Unix(0, 0),
		CommittedAt:        time.Unix(0, 0),
		DeployedAt:         time.Unix(0, 0),
		BuildConfiguration: model.BuildConfiguration{},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		State:              model.DeploymentStateReady,
	}

	expected := &Deployment{
		ID:                 "dpl_test",
		AppID:              "app_test",
		App:                App{},
		Creator:            model.GetCreatorString(model.Actor{}),
		Meta:               "dummy_meta",
		GitSource:          model.GetGitSourceString(model.GitSource{}),
		BuiltAt:            time.Unix(0, 0),
		CommittedAt:        time.Unix(0, 0),
		DeployedAt:         time.Unix(0, 0),
		BuildConfiguration: model.GetBuildConfigurationString(model.BuildConfiguration{}),
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		State:              model.DeploymentStateReady,
	}

	actual := NewDeploymentFromModel(&dpl)
	assert.Equal(t, expected, actual)
}
