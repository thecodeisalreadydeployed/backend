package kubernetesinteractor

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesInteractor struct {
	Clientset *kubernetes.Clientset
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
