package kubernetesinteractor

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesInteractor interface {
	CreatePod(pod apiv1.Pod, namespace string) (string, error)
	CreateConfigMap(configMap apiv1.ConfigMap, namespace string) (string, error)
	GetDeploymentState(name string, namespace string) (model.DeploymentState, error)
}

type kubernetesInteractor struct {
	client *kubernetes.Clientset
}

func NewKubernetesInteractor() (KubernetesInteractor, error) {
	var kubeconfig *string

	if home, err := homedir.Dir(); err != nil {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	// Use the current context in kubeconfig file.
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		zap.L().Error("error creating rest.Config when using BuildConfigFromFlags()")
		return nil, errutil.ErrFailedPrecondition
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		zap.L().Error("error creating kubernetes.Clientset when using NewForConfig()")
		return nil, errutil.ErrFailedPrecondition
	}

	return kubernetesInteractor{client: clientset}, nil
}

func (it kubernetesInteractor) CreatePod(pod apiv1.Pod, namespace string) (string, error) {
	_, err := it.client.CoreV1().Namespaces().Get(context.TODO(), namespace, v1.GetOptions{})
	if err != nil {
		zap.L().Error("namespace not found: " + namespace)
		if _, err := it.createNamespace(namespace); err != nil {
			return "", errutil.ErrFailedPrecondition
		}
	}

	_, err = it.client.CoreV1().Pods(namespace).Get(context.TODO(), pod.Name, v1.GetOptions{})
	if err != nil {
		zap.L().Sugar().Errorf("pod already exists: %s/%s", namespace, pod.Name)
		return "", errutil.ErrAlreadyExists
	}

	create, createErr := it.client.CoreV1().Pods(namespace).Create(context.TODO(), &pod, v1.CreateOptions{})
	if createErr != nil {
		zap.L().Sugar().Errorf("error creating pod: %s/%s", namespace, pod.Name)
		return "", errutil.ErrUnknown

	}

	return create.Name, nil
}

func (it kubernetesInteractor) CreateConfigMap(configMap apiv1.ConfigMap, namespace string) (string, error) {
	_, err := it.client.CoreV1().Namespaces().Get(context.TODO(), namespace, v1.GetOptions{})
	if err != nil {
		zap.L().Error("namespace not found: " + namespace)
		if _, err := it.createNamespace(namespace); err != nil {
			return "", errutil.ErrFailedPrecondition
		}
	}

	_, err = it.client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMap.Name, v1.GetOptions{})
	if err != nil {
		zap.L().Sugar().Errorf("ConfigMap already exists: %s/%s", namespace, configMap.Name)
		return "", errutil.ErrAlreadyExists
	}

	create, createErr := it.client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), &configMap, v1.CreateOptions{})
	if createErr != nil {
		zap.L().Sugar().Errorf("error creating ConfigMap: %s/%s", namespace, &configMap.Name)
		return "", errutil.ErrUnknown

	}

	return create.Name, nil
}

func (it kubernetesInteractor) GetDeploymentState(deploymentID string, namespace string) (model.DeploymentState, error) {
	_, err := it.client.CoreV1().Namespaces().Get(context.TODO(), namespace, v1.GetOptions{})
	if err != nil {
		zap.L().Error("namespace not found: " + namespace)
		return "", errutil.ErrFailedPrecondition
	}

	pods, err := it.client.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{LabelSelector: "codedeploy/deployment-id=" + deploymentID})
	if err != nil {
		zap.L().Sugar().Errorf("error listing pods with codedeploy/deployment-id=%s label in %s namespace", deploymentID, namespace)
		return "", errutil.ErrFailedPrecondition
	}

	if len(pods.Items) == 0 {
		zap.L().Sugar().Errorf("cannot find pod with codedeploy/deployment-id=%s label in %s namespace", deploymentID, namespace)
		return "", errutil.ErrNotFound
	}

	return model.DeploymentStateError, nil
}

func (it kubernetesInteractor) createNamespace(namespace string) (string, error) {
	n := apiv1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: namespace,
		},
	}

	create, createErr := it.client.CoreV1().Namespaces().Create(context.TODO(), &n, v1.CreateOptions{})
	if createErr != nil {
		zap.L().Sugar().Errorf("error creating namespace: %s", namespace)
		return "", errutil.ErrUnknown
	}

	return create.Name, nil
}
