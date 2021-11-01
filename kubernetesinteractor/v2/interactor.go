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
	GetDeploymentState(name string, namespace string) model.DeploymentState
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
		return "", errutil.ErrFailedPrecondition
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

func (kubernetesInteractor) GetDeploymentState(name string, namespace string) model.DeploymentState {
	return model.DeploymentStateError
}
