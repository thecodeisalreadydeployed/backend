package kanikointeractor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor"
)

type KanikoInteractor struct {
	Registry     containerregistry.ContainerRegistryType
	BuildContext string
	Destination  string
}

func (it *KanikoInteractor) BuildContainerImage() error {
	if !strings.HasPrefix(it.BuildContext, "git") {
		return errors.New(fmt.Sprintf("Build context %s is not supported.", it.BuildContext))
	}

	k8s := kubernetesinteractor.NewKubernetesInteractor()

	switch it.Registry {
	case containerregistry.GCR:
		podSpec := it.GCRKanikoPodSpec()
		k8s.CreatePod(&podSpec)
	case containerregistry.ECR:
		podSpec := it.ECRKanikoPodSpec()
		k8s.CreatePod(&podSpec)
	}

	return nil
}
