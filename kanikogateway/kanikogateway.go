package kanikogateway

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

type KanikoGateway struct {
	registry     containerregistry.ContainerRegistry
	buildContext string
	appID        string
	deploymentID string
}

func NewKanikoGateway(registry containerregistry.ContainerRegistry, buildContext string, appID string, deploymentID string) *KanikoGateway {
	return &KanikoGateway{
		registry:     registry,
		buildContext: buildContext,
		appID:        appID,
		deploymentID: deploymentID,
	}
}

func (it *KanikoGateway) BuildContainerImage() error {
	k8s := kubernetesinteractor.NewKubernetesInteractor()

	switch it.registry.Type() {
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

func (it *KanikoGateway) Destination() string {
	dst, _ := it.registry.RegistryFormat(it.appID, it.deploymentID)
	return dst
}

func DeploymentState(deploymentID string) model.DeploymentState {
	k8s := kubernetesinteractor.NewKubernetesInteractor()
	namespace := v1.NamespaceDefault
	podName := fmt.Sprintf("kaniko-%s", deploymentID)
	pod, podsErr := k8s.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})

	if podsErr != nil {
		zap.L().Sugar().Errorf("Cannot find Pod '%s' in namespace '%s'.", podName, namespace)
		return model.DeploymentStateError
	}

	if pod.Status.Phase == v1.PodSucceeded {
		return model.DeploymentStateBuildSucceeded
	}

	if pod.Status.Phase == v1.PodPending || pod.Status.Phase == v1.PodRunning {
		return model.DeploymentStateBuilding
	}

	return model.DeploymentStateError
}
