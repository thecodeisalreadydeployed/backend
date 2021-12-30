package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type Deployment struct {
	ID    string `json:"id"`
	AppID string `json:"appID"`

	Creator Actor  `json:"creator"`
	Meta    string `json:"meta"`

	GitSource GitSource `json:"gitSource"`

	// The time that Deployment.State transitioned to DeploymentStateBuildSucceeded.
	BuiltAt time.Time `json:"builtAt"`

	// The time that Deployment.State transitioned to DeploymentStateCommitted.
	CommittedAt time.Time `json:"committedAt"`

	// The time that Deployment.State transitioned to DeploymentStateReady.
	DeployedAt time.Time `json:"deployedAt"`

	BuildConfiguration BuildConfiguration `json:"buildConfiguration"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	State     DeploymentState `json:"state"`
}

func GenerateDeploymentID() string {
	return fmt.Sprintf("dpl-%s", util.RandomString(25))
}
