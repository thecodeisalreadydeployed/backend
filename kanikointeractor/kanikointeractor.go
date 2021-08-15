package kanikointeractor

import (
	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor"
	apiv1 "k8s.io/api/core/v1"
)

type KanikoInteractor struct {
	Registry containerregistry.ContainerRegistryType
}

func (it *KanikoInteractor) baseKanikoPodSpec() apiv1.Pod {
	podSpec := apiv1.Pod{}
	return podSpec
}

func (it *KanikoInteractor) gcrKanikoPodSpec() apiv1.Pod {
	podSpec := it.baseKanikoPodSpec()
	return podSpec
}

func (it *KanikoInteractor) ecrKanikoPodSpec() apiv1.Pod {
	podSpec := it.baseKanikoPodSpec()
	return podSpec
}

func (it *KanikoInteractor) BuildContainerImage() {
	k8s := kubernetesinteractor.NewKubernetesInteractor()
	k8s.CreatePod(&apiv1.Pod{})
}
