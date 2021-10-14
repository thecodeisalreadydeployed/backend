package manifestgenerator

import (
	"github.com/ghodss/yaml"
	api "sigs.k8s.io/kustomize/api/types"
)

type GenerateKustomizationOptions struct {
	Name      string
	Namespace string
	Resources []string
}

func GenerateKustomizationConfiguration(opts *GenerateKustomizationOptions) (string, error) {
	kustomize := api.Kustomization{
		TypeMeta: api.TypeMeta{
			APIVersion: "kustomize.config.k8s.io/v1beta1",
			Kind:       "Kustomization",
		},
		Resources: opts.Resources,
	}

	y, err := yaml.Marshal(kustomize)

	return string(y), err
}
