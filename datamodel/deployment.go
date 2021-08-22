package datamodel

import (
	"encoding/json"
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type Deployment struct {
	ID        string `gorm:"primaryKey"`
	AppID     string
	App       App
	Name      string
	Creator   string
	Meta      string
	GitSource string
	BuildedAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	State     model.DeploymentState
}

func (dpl *Deployment) ToModel() model.Deployment {
	gitSource := model.GitSource{}
	err := json.Unmarshal([]byte(dpl.GitSource), &gitSource)
	if err != nil {
		panic(err)
	}

	creator := model.Actor{}
	err = json.Unmarshal([]byte(dpl.Creator), &creator)
	if err != nil {
		panic(err)
	}

	return model.Deployment{
		ID:        dpl.ID,
		Name:      dpl.Name,
		Creator:   creator,
		Meta:      dpl.Meta,
		GitSource: gitSource,
		BuildedAt: dpl.BuildedAt,
		CreatedAt: dpl.CreatedAt,
		UpdatedAt: dpl.UpdatedAt,
		State:     dpl.State,
	}
}
