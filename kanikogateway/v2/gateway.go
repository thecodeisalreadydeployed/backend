package kanikogateway

import (
	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor/v2"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

type KanikoGateway interface {
	BuildContainerImage() (string, error)
}

type kanikoGateway struct {
	deploymentID       string
	repositoryURL      string
	branch             string
	buildConfiguration model.BuildConfiguration
	kubernetes         *kubernetesinteractor.KubernetesInteractor
	registry           *containerregistry.ContainerRegistry
}

func NewKanikoGateway(deploymentID string, repositoryURL string, branch string, buildConfiguration model.BuildConfiguration) (KanikoGateway, error) {
	it, err := kubernetesinteractor.NewKubernetesInteractor()

	if err != nil {
		zap.L().Error("could not initialize KubernetesInteractor", zap.String("deploymentID", deploymentID))
		return nil, errutil.ErrFailedPrecondition
	}

	return kanikoGateway{deploymentID: deploymentID, repositoryURL: repositoryURL, branch: branch, buildConfiguration: buildConfiguration, kubernetes: &it, registry: nil}, nil
}

func (kanikoGateway) BuildContainerImage() (string, error) {
	return "", errutil.ErrNotImplemented
}
