package manifestgenerator

import (
	"github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type GenerateServiceOptions struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Selector  map[string]string
}

func GenerateServiceYAML(opts *GenerateServiceOptions) (string, error) {
	srv := apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
			Labels:    opts.Labels,
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			Ports: []apiv1.ServicePort{
				{
					Port:       3000,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: intstr.FromInt(3000),
					Name:       "http",
				},
			},
			Selector: opts.Selector,
		},
	}

	y, err := yaml.Marshal(srv)

	return string(y), err
}
