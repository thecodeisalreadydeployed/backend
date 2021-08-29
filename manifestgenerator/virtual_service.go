package manifestgenerator

import (
	"github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GenerateServiceOptions struct {
	Name           string
	Namespace      string
	Labels         map[string]string
	ContainerImage string
}

func GenerateServiceYAML(opts *GenerateServiceOptions) (string, error) {
	srv := apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
			Labels:    opts.Labels,
		},
		Spec: apiv1.ServiceSpec{},
	}

	y, err := yaml.Marshal(srv)

	return string(y), err
}
