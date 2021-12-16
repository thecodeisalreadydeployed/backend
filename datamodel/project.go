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

func (p *Project) ToModel() model.Project {
	return model.Project{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func NewProjectFromModel(prj *model.Project) *Project {
	return &Project{
		ID:        prj.ID,
		Name:      prj.Name,
		CreatedAt: prj.CreatedAt,
		UpdatedAt: prj.UpdatedAt,
	}
}

func ProjectStructString() []string {
	return []string{"ID", "Name", "CreatedAt", "UpdatedAt"}
}
