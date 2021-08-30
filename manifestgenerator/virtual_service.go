package manifestgenerator

import (
	"github.com/ghodss/yaml"
	nginx "github.com/nginxinc/kubernetes-ingress/pkg/apis/configuration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GenerateVirtualServerOptions struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

func GenerateVirtualServerYAML(opts *GenerateVirtualServerOptions) (string, error) {
	virtualServer := nginx.VirtualServer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.nginx.org/v1",
			Kind:       "VirtualServer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
			Labels:    opts.Labels,
		},
		Spec: nginx.VirtualServerSpec{},
	}

	y, err := yaml.Marshal(virtualServer)

	return string(y), err
}
