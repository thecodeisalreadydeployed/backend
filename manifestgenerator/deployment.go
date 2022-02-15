package manifestgenerator

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GenerateDeploymentOptions struct {
	Name           string
	Namespace      string
	Labels         map[string]string
	Selector       map[string]string
	ContainerImage string
}

func GenerateDeploymentYAML(opts *GenerateDeploymentOptions) (string, error) {
	dpl := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
			Labels:    opts.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: opts.Selector,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: opts.Selector,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            opts.Name,
							Image:           opts.ContainerImage,
							ImagePullPolicy: apiv1.PullIfNotPresent,
							Env: []apiv1.EnvVar{
								{
									Name:  "PORT",
									Value: "3000",
								},
							},
						},
						{
							Name:            "metadataserver",
							Image:           "ghcr.io/thecodeisalreadydeployed/metadataserver:latest",
							ImagePullPolicy: apiv1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}

	y, err := yaml.Marshal(dpl)

	return string(y), err
}
