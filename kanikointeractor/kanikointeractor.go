package kanikointeractor

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func DeploymentState(deploymentID string) model.DeploymentState {
	k8s := kubernetesinteractor.NewKubernetesInteractor()
	namespace := v1.NamespaceDefault
	podName := fmt.Sprintf("kaniko-%s", deploymentID)
	pod, podsErr := k8s.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})

	if podsErr != nil {
		zap.L().Sugar().Errorf("Cannot find pod '%s' in namespace '%s'.", podName, namespace)
		return model.DeploymentStateError
	}

	if pod.Status.Phase == v1.PodSucceeded {
		// TODO: Replace with model.DeploymentStateBuildSucceeded
		return model.DeploymentStateBuilding
	}

	return model.DeploymentStateError
}
