package kanikointeractor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KanikoInteractor struct {
	Registry     containerregistry.ContainerRegistryType
	BuildContext string
}

func (it *KanikoInteractor) baseKanikoPodSpec() apiv1.Pod {
	podSpec := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kaniko",
		},
		Spec: apiv1.PodSpec{
			RestartPolicy: apiv1.RestartPolicyNever,
			Containers: []apiv1.Container{
				{
					Name:  "kaniko",
					Image: "gcr.io/kaniko-project/executor:v1.6.0",
				},
			},
		},
	}
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

func (it *KanikoInteractor) BuildContainerImage() error {
	if !strings.HasPrefix(it.BuildContext, "git") {
		return errors.New(fmt.Sprintf("Build context %s is not supported.", it.BuildContext))
	}

	k8s := kubernetesinteractor.NewKubernetesInteractor()
	k8s.CreatePod(&apiv1.Pod{})

	return nil
}
