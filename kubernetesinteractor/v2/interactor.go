package kubernetesinteractor

import (
	"flag"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
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

func (kubernetesInteractor) CreatePod(pod apiv1.Pod, namespace string) (string, error) {
	return "", errutil.ErrNotImplemented
}

func (kubernetesInteractor) GetDeploymentState(name string, namespace string) model.DeploymentState {
	return model.DeploymentStateError
}
