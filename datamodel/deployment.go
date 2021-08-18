package datamodel

import (
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type DeploymentState string

const (
	DeploymentStateQueueing DeploymentState = "DeploymentStateQueueing"
	DeploymentStateBuilding DeploymentState = "DeploymentStateBuilding"
	DeploymentStateReady    DeploymentState = "DeploymentStateReady"
	DeploymentStateError    DeploymentState = "DeploymentStateError"
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
	State     DeploymentState
}

func (dpl *Deployment) toModel() model.Deployment {}
