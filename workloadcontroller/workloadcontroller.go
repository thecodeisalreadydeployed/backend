package workloadcontroller

import (
	"fmt"

	"github.com/thecodeisalreadydeployed/containerregistry/gcr"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/kanikogateway"
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

func NewDeployment(appID string) (string, error) {
	app, getAppErr := datastore.GetAppByID(appID)
	if getAppErr != nil {
		return "", getAppErr
	}

	buildContext := app.GitSource.RepositoryURL

	deploymentID := util.RandomString(5)
	containerRegistry := gcr.NewGCRGateway("asia.gcr.io", "hu-tao-mains", "")

	kaniko := kanikogateway.NewKanikoGateway(containerRegistry, buildContext, app.ID, deploymentID)

	buildContainerImageErr := kaniko.BuildContainerImage()
	if buildContainerImageErr != nil {
		return "", buildContainerImageErr
	}

	return deploymentID, nil
}

func OnGitSourceUpdate(shouldDeploy bool) {}

func DeployNewRevision() {
	fmt.Println("New revision deploying...")
}
