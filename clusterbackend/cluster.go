package clusterbackend

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ClusterBackend interface {
	CreateNamespace(name string) (string, error)
	CreatePod(pod apiv1.Pod, namespace string) (string, error)
	CreateConfigMap(configMap apiv1.ConfigMap, namespace string) (string, error)
}

type clusterBackend struct {
	logger           *zap.Logger
	kubernetesClient *kubernetes.Clientset
}

func NewClusterBackend(logger *zap.Logger) ClusterBackend {
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Fatal("cannot create in-cluster config", zap.Error(err))
	}

	kubernetesClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal("cannot create Kubernetes client from in-cluster config", zap.Error(err))
	}

	backend := &clusterBackend{logger: logger, kubernetesClient: kubernetesClient}
	return backend
}

func (backend *clusterBackend) CreateNamespace(name string) (string, error) {
	n := apiv1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
	}

	createdNamespace, createErr := backend.kubernetesClient.CoreV1().Namespaces().Create(context.TODO(), &n, v1.CreateOptions{})
	if createErr != nil {
		backend.logger.Error("cannot create namespace", zap.String("namespace", name), zap.Error(createErr))
		return "", fmt.Errorf("cannot create namespace %s: %w", name, createErr)
	}

	return createdNamespace.Name, nil
}

func (backend *clusterBackend) CreatePod(pod apiv1.Pod, namespace string) (string, error) {
	return "", nil
}

func (backend *clusterBackend) CreateConfigMap(configMap apiv1.ConfigMap, namespace string) (string, error) {
	return "", nil
}
