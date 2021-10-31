package kanikogateway

import (
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
)

type KanikoGateway interface {
	BuildContainerImage() (string, error)
}

type kanikoGateway struct {
	deploymentID       string
	repositoryURL      string
	branch             string
	buildConfiguration model.BuildConfiguration
}

func NewKanikoGateway(deploymentID string, repositoryURL string, branch string, buildConfiguration model.BuildConfiguration) KanikoGateway {
	return kanikoGateway{deploymentID: deploymentID, repositoryURL: repositoryURL, branch: branch, buildConfiguration: buildConfiguration}
}

func (kanikoGateway) BuildContainerImage() (string, error) {
	return "", errutil.ErrNotImplemented
}
