package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type Deployment struct {
	ID    string `json:"id"`
	AppID string `json:"app_id"`

	Creator Actor  `json:"creator"`
	Meta    string `json:"meta"`

	GitSource GitSource `json:"git_source"`

	// The time that Deployment.State transitioned to DeploymentStateBuildSucceeded.
	BuiltAt time.Time `json:"built_at"`

	// The time that Deployment.State transitioned to DeploymentStateCommitted.
	CommittedAt time.Time `json:"committed_at"`

	// The time that Deployment.State transitioned to DeploymentStateReady.
	DeployedAt time.Time `json:"deployed_at"`

	// The Dockerfile instructions that is used to build the container image.
	BuildScript string `json:"build_script"`

	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	State     DeploymentState `json:"state"`
}

func GenerateDeploymentID() string {
	return fmt.Sprintf("dpl_%s", util.RandomString(5))
}
