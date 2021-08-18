package datamodel

import (
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type App struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	GitSource string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (app *App) toModel() model.App {
	return model.App{
		ID:        app.ID,
		Name:      app.Name,
		GitSource: app.GitSource,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}
}
