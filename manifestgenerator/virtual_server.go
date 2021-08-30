package manifestgenerator

import (
	"fmt"

	"github.com/ghodss/yaml"
	nginx "github.com/nginxinc/kubernetes-ingress/pkg/apis/configuration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GenerateVirtualServerOptions struct {
	Labels map[string]string

	ProjectID string
	AppID     string
}

func GenerateVirtualServerYAML(opts *GenerateVirtualServerOptions) (string, error) {
	virtualServer := nginx.VirtualServer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.nginx.org/v1",
			Kind:       "VirtualServer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.AppID,
			Namespace: opts.ProjectID,
			Labels:    opts.Labels,
		},
		Spec: nginx.VirtualServerSpec{
			Host: fmt.Sprintf("%s.%s", opts.AppID, ""),
			TLS: &nginx.TLS{
				Secret: "",
			},
			Upstreams: []nginx.Upstream{
				{
					Name:    opts.AppID,
					Service: opts.AppID,
					Port:    uint16(3000),
				},
			},
			Routes: []nginx.Route{
				{
					Path: "/",
					Action: &nginx.Action{
						Pass: opts.AppID,
					},
				},
			},
		},
	}

	y, err := yaml.Marshal(virtualServer)

	return string(y), err
}
