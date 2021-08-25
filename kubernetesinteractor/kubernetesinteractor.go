package kubernetesinteractor

import (
	"context"
	"fmt"

	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesInteractor struct {
	Clientset *kubernetes.Clientset
}

func NewKubernetesInteractor() KubernetesInteractor {
	var kubeconfig string
	pflag.StringVar(&kubeconfig, "kubeconfig", "", "")
	pflag.Parse()
	kubeconfig = absolutePath(kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return KubernetesInteractor{
		Clientset: clientset,
	}
}

func (it *KubernetesInteractor) ListDeployments() {
	deploymentsClient := it.Clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployments, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	spew.Dump(deployments)
}

func (it *KubernetesInteractor) CreateDeployment(deployment *appsv1.Deployment) {
	deploymentsClient := it.Clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	spew.Dump(result.GetObjectMeta())
}

func (it *KubernetesInteractor) CreatePod(pod *apiv1.Pod) {
	podsClient := it.Clientset.CoreV1().Pods(apiv1.NamespaceDefault)
	result, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	spew.Dump(result.GetObjectMeta())
}

func absolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	if strings.HasPrefix(path, "~") {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		return filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	return path
}
