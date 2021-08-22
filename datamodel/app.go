package datamodel

import (
	"encoding/json"
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type App struct {
	ID              string `gorm:"primaryKey"`
	ProjectID       string
	Project         Project
	Name            string
	GitSource       string
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	BuildCommand    string
	OutputDirectory string
	InstallCommand  string
}

type BareApp struct {
	ID              string
	ProjectID       string
	Name            string
	GitSource       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	BuildCommand    string
	OutputDirectory string
	InstallCommand  string
}

func (app *App) toModel() model.App {
	gitSource := model.GitSource{}
	err := json.Unmarshal([]byte(app.GitSource), &gitSource)
	if err != nil {
		panic(err)
	}
	return model.App{
		ID:        app.ID,
		Name:      app.Name,
		GitSource: gitSource,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
		BuildCommand:    app.BuildCommand,
		OutputDirectory: app.OutputDirectory,
		InstallCommand:  app.InstallCommand,
	}
}
