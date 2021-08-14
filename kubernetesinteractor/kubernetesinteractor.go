package kubernetesinteractor

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListDeployments(clientset *kubernetes.Clientset) {
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployments, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	spew.Dump(deployments)
}
