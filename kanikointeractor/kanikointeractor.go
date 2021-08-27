package kanikointeractor

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor"
)

type KanikoInteractor struct {
	Registry     containerregistry.ContainerRegistry
	BuildContext string
	Destination  string
	DeploymentID string
}

func (it *KanikoInteractor) BuildContainerImage() error {
	k8s := kubernetesinteractor.NewKubernetesInteractor()

	switch it.Registry.Type() {
	case containerregistry.GCR:
		podSpec := it.GCRKanikoPodSpec()
		spew.Dump(podSpec)
		k8s.CreatePod(&podSpec)
	case containerregistry.ECR:
		podSpec := it.ECRKanikoPodSpec()
		k8s.CreatePod(&podSpec)
	case containerregistry.LOCAL:
		podSpec := it.baseKanikoPodSpec()
		k8s.CreatePod(&podSpec)
	}

	return nil
}
