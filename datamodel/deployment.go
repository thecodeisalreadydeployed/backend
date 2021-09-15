package datamodel

import (
	"encoding/json"
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type Deployment struct {
	ID          string `gorm:"primaryKey"`
	AppID       string
	App         App `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Creator     string
	Meta        string
	GitSource   string
	BuiltAt     time.Time
	CommittedAt time.Time
	DeployedAt  time.Time
	BuildScript string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	State       model.DeploymentState
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
		ID:          dpl.ID,
		Creator:     creator,
		Meta:        dpl.Meta,
		GitSource:   gitSource,
		BuiltAt:     dpl.BuiltAt,
		CommittedAt: dpl.CommittedAt,
		DeployedAt:  dpl.DeployedAt,
		BuildScript: dpl.BuildScript,
		CreatedAt:   dpl.CreatedAt,
		UpdatedAt:   dpl.UpdatedAt,
		State:       dpl.State,
	}
}
