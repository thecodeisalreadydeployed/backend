package kanikointeractor

import (
	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor"
	v1 "k8s.io/api/apps/v1"
)

type KanikoInteractor struct {
	Registry containerregistry.ContainerRegistryType
}

func (it KanikoInteractor) BuildContainerImage() {
	k8s := kubernetesinteractor.NewKubernetesInteractor()
	k8s.CreateDeployment(&v1.Deployment{})
}
