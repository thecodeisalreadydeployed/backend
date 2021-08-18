package datamodel

import (
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type Project struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (p *Project) toModel() model.Project {
	return model.Project{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
