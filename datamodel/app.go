package datamodel

import (
	"encoding/json"
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type App struct {
	ID              string `gorm:"primaryKey"`
	ProjectID       string
	Project         Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name            string
	GitSource       string
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	BuildCommand    string
	OutputDirectory string
	InstallCommand  string
}

func (app *App) ToModel() model.App {
	gitSource := model.GitSource{}
	err := json.Unmarshal([]byte(app.GitSource), &gitSource)
	if err != nil {
		panic(err)
	}
	return model.App{
		ID:              app.ID,
		Name:            app.Name,
		GitSource:       gitSource,
		CreatedAt:       app.CreatedAt,
		UpdatedAt:       app.UpdatedAt,
		BuildCommand:    app.BuildCommand,
		OutputDirectory: app.OutputDirectory,
		InstallCommand:  app.InstallCommand,
	}
}
