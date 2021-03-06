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
							Name:            "container0",
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
							Env: []apiv1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								{
									Name: "POD_NAME",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name: "PROJECT_ID",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.labels['beta.deploys.dev/project-id']",
										},
									},
								},
								{
									Name: "APP_ID",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.labels['beta.deploys.dev/app-id']",
										},
									},
								},
								{
									Name: "DEPLOYMENT_ID",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.labels['beta.deploys.dev/deployment-id']",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	y, err := yaml.Marshal(dpl)

	return string(y), err
}
