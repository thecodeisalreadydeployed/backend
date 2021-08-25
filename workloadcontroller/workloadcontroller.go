package workloadcontroller

import (
	apidto "github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/manifestgenerator"
	manifestdto "github.com/thecodeisalreadydeployed/manifestgenerator/dto"
)

func CreateWorkload(w *apidto.CreateProjectRequest) string {
	spec := manifestdto.ContainerSpec{
		Name:  "nginx-container",
		Image: "nginx",
		Port:  8000,
	}

	dpl := manifestdto.Deployment{
		APIVersion:    "v1",
		Name:          "test-deploy",
		Replicas:      3,
		Labels:        map[string]string{"app": "nginx"},
		ContainerSpec: spec,
	}

	return manifestgenerator.CreateDeploymentYAML(&dpl)
}
