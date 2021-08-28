package workloadcontroller

import (
	manifest "github.com/thecodeisalreadydeployed/manifestgenerator"
	"github.com/thecodeisalreadydeployed/util"
)

type NewAppOptions struct{}

func NewApp(opts *NewAppOptions) bool {
	manifest.GenerateDeploymentYAML(&manifest.GenerateDeploymentOptions{
		Name:           util.RandomString(5),
		Namespace:      util.RandomString(5),
		Labels:         map[string]string{},
		ContainerImage: "k8s.gcr.io/ingress-nginx/controller:v1.0.0",
	})

	return true
}
