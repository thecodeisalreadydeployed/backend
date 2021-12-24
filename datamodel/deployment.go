package datamodel

import (
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type Deployment struct {
	ID                 string `gorm:"primaryKey"`
	AppID              string
	App                App `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Creator            string
	Meta               string
	GitSource          string
	BuiltAt            time.Time
	CommittedAt        time.Time
	DeployedAt         time.Time
	BuildConfiguration string
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
	State              model.DeploymentState
}

func (dpl *Deployment) ToModel() model.Deployment {
	return model.Deployment{
		ID:                 dpl.ID,
		AppID:              dpl.AppID,
		Creator:            model.GetCreatorObject(dpl.Creator),
		Meta:               dpl.Meta,
		GitSource:          model.GetGitSourceObject(dpl.GitSource),
		BuiltAt:            dpl.BuiltAt,
		CommittedAt:        dpl.CommittedAt,
		DeployedAt:         dpl.DeployedAt,
		BuildConfiguration: model.GetBuildConfigurationObject(dpl.BuildConfiguration),
		CreatedAt:          dpl.CreatedAt,
		UpdatedAt:          dpl.UpdatedAt,
		State:              dpl.State,
	}
}

func NewDeploymentFromModel(dpl *model.Deployment) *Deployment {
	return &Deployment{
		ID:                 dpl.ID,
		AppID:              dpl.AppID,
		Meta:               dpl.Meta,
		State:              dpl.State,
		GitSource:          model.GetGitSourceString(dpl.GitSource),
		Creator:            model.GetCreatorString(dpl.Creator),
		BuildConfiguration: model.GetBuildConfigurationString(dpl.BuildConfiguration),
		BuiltAt:            dpl.BuiltAt,
		CommittedAt:        dpl.CommittedAt,
		DeployedAt:         dpl.DeployedAt,
		CreatedAt:          dpl.CreatedAt,
		UpdatedAt:          dpl.UpdatedAt,
	}
}

func DeploymentStructString() []string {
	return []string{
		"ID",
		"AppID",
		"Creator",
		"Meta",
		"GitSource",
		"BuiltAt",
		"CommittedAt",
		"DeployedAt",
		"BuildConfiguration",
		"CreatedAt",
		"UpdatedAt",
		"State",
	}
}
