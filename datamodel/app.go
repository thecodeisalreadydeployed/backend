package datamodel

import (
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type App struct {
	ID                 string `gorm:"primaryKey"`
	ProjectID          string
	Project            Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name               string
	GitSource          string
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
	BuildConfiguration string
	Observable         bool
}

func (app *App) ToModel() model.App {
	return model.App{
		ID:                 app.ID,
		ProjectID:          app.ProjectID,
		Name:               app.Name,
		GitSource:          model.GetGitSourceObject(app.GitSource),
		CreatedAt:          app.CreatedAt,
		UpdatedAt:          app.UpdatedAt,
		BuildConfiguration: model.GetBuildConfigurationObject(app.BuildConfiguration),
		Observable:         app.Observable,
	}
}

func NewAppFromModel(app *model.App) *App {
	return &App{
		ID:                 app.ID,
		ProjectID:          app.ProjectID,
		Name:               app.Name,
		GitSource:          model.GetGitSourceString(app.GitSource),
		CreatedAt:          app.CreatedAt,
		UpdatedAt:          app.UpdatedAt,
		BuildConfiguration: model.GetBuildConfigurationString(app.BuildConfiguration),
		Observable:         app.Observable,
	}
}
