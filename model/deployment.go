package model

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/util"
)

type Deployment struct {
	ID        string          `json:"id"`
	AppID     string          `json:"app_id"`
	Name      string          `json:"name"`
	Creator   Actor           `json:"creator"`
	Meta      string          `json:"meta"`
	GitSource GitSource       `json:"git_source"`
	BuiltAt   time.Time       `json:"builded_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	State     DeploymentState `json:"state"`
}

func GenerateDeploymentID() string {
	return fmt.Sprintf("dpl_%s", util.RandomString(5))
}
