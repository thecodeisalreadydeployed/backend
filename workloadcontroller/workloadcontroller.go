package workloadcontroller

import (
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
