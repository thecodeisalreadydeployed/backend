package manifestgenerator

import (
	"github.com/thecodeisalreadydeployed/manifestgenerator/dto"
	"github.com/thecodeisalreadydeployed/manifestgenerator/json"
	"gopkg.in/yaml.v2"
)

func CreateDeploymentYAML(dpl *dto.Deployment) string {
	var port json.Port
	port.ContainerPort = dpl.ContainerSpec.Port

	var ctn json.Container
	ctn.ContainerName = dpl.ContainerSpec.Name
	ctn.ContainerImage = dpl.ContainerSpec.Image
	ctn.ContainerPorts = []json.Port{port}

	var obj json.Deployment
	obj.ApiVersion = dpl.APIVersion
	obj.Kind = "Deployment"
	obj.Metadata.Name = dpl.Name
	obj.Metadata.Labels = dpl.Labels
	obj.Spec.Replicas = dpl.Replicas
	obj.Spec.Selector.Labels = dpl.Labels
	obj.Spec.Template.TemplateMetadata.Labels = dpl.Labels
	obj.Spec.Template.TemplateSpec.Containers = []json.Container{ctn}

	y, err := yaml.Marshal(&obj)
	if err != nil {
		panic(err)
	}
	return string(y)
}
