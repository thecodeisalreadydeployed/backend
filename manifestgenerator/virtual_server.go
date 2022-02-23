package manifestgenerator

import (
	"fmt"

	"github.com/ghodss/yaml"
	nginx "github.com/thecodeisalreadydeployed/manifestgenerator/nginx"
)

type GenerateVirtualServerOptions struct {
	Labels map[string]string

	ProjectID string
	AppID     string
}

func GenerateVirtualServerYAML(opts *GenerateVirtualServerOptions) (string, error) {
	virtualServer := nginx.VirtualServer{
		APIVersion: "k8s.nginx.org/v1",
		Kind:       "VirtualServer",
		Metadata: nginx.Metadata{
			Name:      opts.AppID,
			Namespace: opts.ProjectID,
			Labels:    opts.Labels,
		},
		Spec: nginx.Spec{
			Host: fmt.Sprintf("%s.svc.deploys.dev", opts.AppID),
			TLS: nginx.TLS{
				Secret: "",
				Redirect: nginx.TLSRedirect{
					Enable: true,
				},
			},
			Upstreams: []nginx.Upstream{
				{
					Name:    opts.AppID,
					Service: opts.AppID,
					Port:    3000,
				},
			},
			Routes: []nginx.Route{
				{
					Path: "/",
					Action: nginx.RouteAction{
						Pass: opts.AppID,
					},
				},
			},
		},
	}

	y, err := yaml.Marshal(virtualServer)

	return string(y), err
}
