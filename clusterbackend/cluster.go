package clusterbackend

import (
	"context"
	"fmt"

	"github.com/thecodeisalreadydeployed/util"
	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ClusterBackend interface {
	CreateNamespace(name string) (string, error)
	CreateNamespaceIfNotExists(name string) (string, error)
	CreatePod(pod apiv1.Pod, namespace string) (string, error)
	CreateConfigMap(configMap apiv1.ConfigMap, namespace string) (string, error)

	Pods(namespace string, labels map[string]string) ([]apiv1.Pod, error)
}

type clusterBackend struct {
	logger           *zap.Logger
	kubernetesClient *kubernetes.Clientset
}

func NewClusterBackend(logger *zap.Logger) ClusterBackend {
	if util.IsDevEnvironment() || util.IsDockerTestEnvironment() {
		return &clusterBackend{logger: logger, kubernetesClient: nil}
	}

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

func (backend *clusterBackend) CreateNamespaceIfNotExists(name string) (string, error) {
	namespace, err := backend.kubernetesClient.CoreV1().Namespaces().Get(context.TODO(), name, v1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		return backend.CreateNamespace(name)
	}

	if err != nil {
		return "", fmt.Errorf("cannot find namespace %s: %w", name, err)
	}

	return namespace.Name, nil
}

func (backend *clusterBackend) CreatePod(pod apiv1.Pod, namespace string) (string, error) {
	_, err := backend.CreateNamespaceIfNotExists(namespace)
	if err != nil {
		return "", err
	}

	_, err = backend.kubernetesClient.CoreV1().Pods(namespace).Get(context.TODO(), pod.Name, v1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		createdPod, createErr := backend.kubernetesClient.CoreV1().Pods(namespace).Create(context.TODO(), &pod, v1.CreateOptions{})
		if createErr != nil {
			backend.logger.Error("cannot create Pod", zap.String("namespace", namespace), zap.String("pod", pod.Name), zap.Error(createErr))
			return "", fmt.Errorf("cannot create Pod %s: %w", pod.Name, createErr)
		}

		return createdPod.Name, nil
	}

	if err != nil {
		backend.logger.Error("cannot create Pod", zap.String("namespace", namespace), zap.String("pod", pod.Name), zap.Error(err))
		return "", fmt.Errorf("cannot create Pod %s: %w", pod.Name, err)
	}

	return "", nil
}

func (backend *clusterBackend) CreateConfigMap(configMap apiv1.ConfigMap, namespace string) (string, error) {
	_, err := backend.CreateNamespaceIfNotExists(namespace)
	if err != nil {
		return "", err
	}

	_, err = backend.kubernetesClient.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMap.Name, v1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		createdConfigMap, createErr := backend.kubernetesClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), &configMap, v1.CreateOptions{})
		if createErr != nil {
			backend.logger.Error("cannot create ConfigMap", zap.String("namespace", namespace), zap.String("configMap", configMap.Name), zap.Error(createErr))
			return "", fmt.Errorf("cannot create ConfigMap %s: %w", configMap.Name, createErr)
		}

		return createdConfigMap.Name, nil
	}

	if err != nil {
		backend.logger.Error("cannot create ConfigMap", zap.String("namespace", namespace), zap.String("configMap", configMap.Name), zap.Error(err))
		return "", fmt.Errorf("cannot create ConfigMap %s: %w", configMap.Name, err)
	}

	return "", nil
}

func (backend *clusterBackend) Pods(namespace string, labels map[string]string) ([]apiv1.Pod, error) {
	labelSelector := v1.LabelSelector{}
	for k, v := range labels {
		v1.AddLabelToSelector(&labelSelector, k, v)
	}

	podList, err := backend.kubernetesClient.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{LabelSelector: v1.FormatLabelSelector(&labelSelector)})
	if err != nil {
		backend.logger.Error("cannot list Pods", zap.String("namespace", namespace), zap.Any("labels", labels), zap.Error(err))
		return []apiv1.Pod{}, fmt.Errorf("cannot list Pods: %w", err)
	}

	return podList.Items, nil
}

func (backend *clusterBackend) DeletePod(namespace string, name string) error {
	_, err := backend.kubernetesClient.CoreV1().Pods(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		return fmt.Errorf("cannot find Pod %s: %w", name, err)
	}

	err = backend.kubernetesClient.CoreV1().Pods(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("cannot delete Pod %s: %w", name, err)
	}

	return nil
}
