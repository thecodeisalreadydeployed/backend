package kubernetesinteractor

import (
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	apiv1 "k8s.io/api/core/v1"
)

type KubernetesInteractor interface {
	CreatePod(pod apiv1.Pod, namespace string) (string, error)
	GetDeploymentState(name string, namespace string) model.DeploymentState
}

type kubernetesInteractor struct{}

func NewKubernetesInteractor() KubernetesInteractor {
	return kubernetesInteractor{}
}

func (kubernetesInteractor) CreatePod(pod apiv1.Pod, namespace string) (string, error) {
	return "", errutil.ErrNotImplemented
}

func (kubernetesInteractor) GetDeploymentState(name string, namespace string) model.DeploymentState {
	return model.DeploymentStateError
}
