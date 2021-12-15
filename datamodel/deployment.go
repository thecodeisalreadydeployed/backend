package datamodel

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/spf13/cast"
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

	buildConfiguration := model.BuildConfiguration{}
	_buildConfiguration, err := base64.StdEncoding.DecodeString(dpl.BuildConfiguration)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(_buildConfiguration, &buildConfiguration)
	if err != nil {
		panic(err)
	}

	return model.Deployment{
		ID:                 dpl.ID,
		AppID:              dpl.AppID,
		Creator:            creator,
		Meta:               dpl.Meta,
		GitSource:          gitSource,
		BuiltAt:            dpl.BuiltAt,
		CommittedAt:        dpl.CommittedAt,
		DeployedAt:         dpl.DeployedAt,
		BuildConfiguration: buildConfiguration,
		CreatedAt:          dpl.CreatedAt,
		UpdatedAt:          dpl.UpdatedAt,
		State:              dpl.State,
	}
}

func NewDeploymentFromModel(dpl *model.Deployment) *Deployment {
	gitSource, err := json.Marshal(dpl.GitSource)
	if err != nil {
		panic(err)
	}

	creator, err := json.Marshal(dpl.Creator)
	if err != nil {
		panic(err)
	}

	buildConfiguration, err := json.Marshal(dpl.BuildConfiguration)
	if err != nil {
		panic(err)
	}

	buildConfiguration64 := base64.StdEncoding.EncodeToString(buildConfiguration)

	return &Deployment{
		ID:                 dpl.ID,
		AppID:              dpl.AppID,
		Meta:               dpl.Meta,
		State:              dpl.State,
		GitSource:          cast.ToString(gitSource),
		Creator:            cast.ToString(creator),
		BuildConfiguration: buildConfiguration64,
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
