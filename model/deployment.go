package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type Deployment struct {
	ID        string
	Name      string
	Creator   Actor
	Meta      string
	GitSource GitSource
	BuildedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	State     DeploymentState
}

func GenerateDeploymentID() string {
	return fmt.Sprintf("dpl_%s", util.RandomString(5))
}
