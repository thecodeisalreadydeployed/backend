package clusterbackend

import (
	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type ClusterBackend interface {
	CreatePod(pod apiv1.Pod, namespace string) (string, error)
	CreateConfigMap(configMap apiv1.ConfigMap, namespace string) (string, error)
}

type clusterBackend struct {
	logger           *zap.Logger
	kubernetesClient *kubernetes.Clientset
}

func NewClusterBackend() ClusterBackend {
	backend := &clusterBackend{}
	return backend
}

func (backend *clusterBackend) CreatePod(pod apiv1.Pod, namespace string) (string, error) {
	return "", nil
}

func (backend *clusterBackend) CreateConfigMap(configMap apiv1.ConfigMap, namespace string) (string, error) {
	return "", nil
}
