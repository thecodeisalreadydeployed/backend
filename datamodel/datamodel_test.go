package datamodel

import (
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/model"
	"testing"
	"time"
)

func TestProject_ToModel(t *testing.T) {
	prj := Project{
		ID:        "prj-test",
		Name:      "Best Project",
		CreatedAt: time.Unix(1, 0),
		UpdatedAt: time.Unix(2, 0),
	}

	parsed := prj.ToModel()
	actual := *NewProjectFromModel(&parsed)
	assert.Equal(t, prj, actual)
}

func TestApp_ToModel(t *testing.T) {
	app := App{
		ID:                 "app-test",
		ProjectID:          "prj-test",
		Project:            Project{},
		Name:               "Best App",
		GitSource:          model.GetGitSourceString(model.GitSource{}),
		CreatedAt:          time.Unix(1, 0),
		UpdatedAt:          time.Unix(2, 0),
		BuildConfiguration: model.GetBuildConfigurationString(model.BuildConfiguration{}),
		Observable:         false,
	}

	parsed := app.ToModel()
	actual := *NewAppFromModel(&parsed)
	assert.Equal(t, app, actual)
}

func TestDeployment_ToModel(t *testing.T) {
	dpl := Deployment{
		ID:                 "dpl-test",
		AppID:              "app-test",
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

	parsed := dpl.ToModel()
	actual := *NewDeploymentFromModel(&parsed)
	assert.Equal(t, dpl, actual)
}

func TestEvent_ToModel(t *testing.T) {
	event := Event{
		ID:           "abcdefghijklmnopqrstuvwxyz0",
		DeploymentID: "dpl-test",
		Text:         "Downloading dependencies (1/20)",
		Type:         model.INFO,
		ExportedAt:   time.Unix(0, 0),
		CreatedAt:    time.Unix(0, 0),
	}

	parsed := event.ToModel()
	actual := *NewEventFromModel(&parsed)
	assert.Equal(t, event, actual)
}

func TestPreset_ToModel(t *testing.T) {
	preset := Preset{
		ID:       "pst-test",
		Name:     "My Preset",
		Template: "UlVOIGVjaG8gaGVsbG8=",
	}

	parsed := preset.ToModel()
	assert.Equal(t, "RUN echo hello", parsed.Template)
	actual := *NewPresetFromModel(&parsed)
	assert.Equal(t, preset, actual)
}
