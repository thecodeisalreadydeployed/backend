package datamodel

import (
	"encoding/base64"
	"encoding/json"
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
}

func (app *App) ToModel() model.App {
	gitSource := model.GitSource{}
	err := json.Unmarshal([]byte(app.GitSource), &gitSource)
	if err != nil {
		panic(err)
	}

	buildConfiguration := model.BuildConfiguration{}
	_buildConfiguration, err := base64.StdEncoding.DecodeString(app.BuildConfiguration)
	err = json.Unmarshal(_buildConfiguration, &buildConfiguration)
	if err != nil {
		panic(err)
	}

	return model.App{
		ID:                 app.ID,
		Name:               app.Name,
		GitSource:          gitSource,
		CreatedAt:          app.CreatedAt,
		UpdatedAt:          app.UpdatedAt,
		BuildConfiguration: buildConfiguration,
	}
}
