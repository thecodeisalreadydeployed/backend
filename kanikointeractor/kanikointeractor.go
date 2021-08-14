package kanikointeractor

import (
	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor"
	apiv1 "k8s.io/api/core/v1"
)

type KanikoInteractor struct {
	Registry containerregistry.ContainerRegistryType
}

func (it KanikoInteractor) BuildContainerImage() {
	k8s := kubernetesinteractor.NewKubernetesInteractor()
	k8s.CreatePod(&apiv1.Pod{})
}
