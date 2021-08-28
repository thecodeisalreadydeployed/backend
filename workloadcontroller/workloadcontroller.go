package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/containerregistry/gcr"
	"github.com/thecodeisalreadydeployed/kanikointeractor"
	manifest "github.com/thecodeisalreadydeployed/manifestgenerator"
	"github.com/thecodeisalreadydeployed/util"
	"go.uber.org/zap"
)

type NewAppOptions struct{}

func NewApp(opts *NewAppOptions) error {
	_, generateDeploymentErr := manifest.GenerateDeploymentYAML(&manifest.GenerateDeploymentOptions{
		Name:           util.RandomString(5),
		Namespace:      util.RandomString(5),
		Labels:         map[string]string{},
		ContainerImage: "k8s.gcr.io/ingress-nginx/controller:v1.0.0",
	})

	if generateDeploymentErr != nil {
		zap.L().Error(generateDeploymentErr.Error())
		return generateDeploymentErr
	}

	return nil
}

type NewDeploymentOptions struct {
	GitRepositoryURL string
}

func NewDeployment(opts *NewDeploymentOptions) (string, error) {
	deploymentID := util.RandomString(5)
	containerRegistry := gcr.NewGCRGateway("asia.gcr.io", "hu-tao-mains", "")
	destination, err := containerRegistry.RegistryFormat("fixture-monorepo", "dev")

	if err != nil {
		return "", err
	}

	// TODO: Rename KanikoInteractor to KanikoGateway
	kaniko := kanikointeractor.KanikoInteractor{
		Registry:     containerRegistry,
		Destination:  destination,
		BuildContext: opts.GitRepositoryURL,
		DeploymentID: deploymentID,
	}

	buildContainerImageErr := kaniko.BuildContainerImage()
	if buildContainerImageErr != nil {
		return "", buildContainerImageErr
	}

	return deploymentID, nil
}
