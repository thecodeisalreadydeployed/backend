package manifestgenerator

import (
	"github.com/thecodeisalreadydeployed/manifestgenerator/dto"
)

func createDeploymentYAML(dpl *dto.Deployment) string {
	var yb YAMLBuilder
	yb.AppendApiVersion(dpl.ApiVersion)
	yb.AppendKind("Deployment")
	yb.AppendDeploymentMetadata(dpl.Name, dpl.Labels)
	yb.EnterSpec()
	yb.AppendReplicas(dpl.Replicas)
	yb.Enter("selector:")
	yb.Enter("matchLabels:")
	yb.AppendLabels(dpl.Labels)
	yb.Close(2)
	yb.EnterTemplate()
	yb.EnterLabels("metadata:", "labels:")
	yb.EnterSpec()
	yb.AppendContainers(dpl.ContainerSpec)
}
