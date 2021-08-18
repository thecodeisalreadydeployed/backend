package datamodel

import (
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type Deployment struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Creator   string
	Meta      string
	GitBranch string
	BuildedAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	State     model.DeploymentState
}

func (dpl *Deployment) toModel() model.Deployment {
	return model.Deployment{
		ID:        dpl.ID,
		Name:      dpl.Name,
		Creator:   dpl.Creator,
		Meta:      dpl.Meta,
		GitBranch: dpl.GitBranch,
		BuildedAt: dpl.BuildedAt,
		CreatedAt: dpl.CreatedAt,
		UpdatedAt: dpl.UpdatedAt,
		State:     dpl.State,
	}
}
