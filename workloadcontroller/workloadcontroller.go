package workloadcontroller

import (
	"fmt"
	apidto "github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/manifestgenerator"
	manifestdto "github.com/thecodeisalreadydeployed/manifestgenerator/dto"
)

func CreateWorkload(w *apidto.CreateProjectRequest) {
	fmt.Printf("Workload %s created.", w.Name)

	spec := manifestdto.ContainerSpec{
		Name: "nginx-container",
		Image: "nginx-fa6ajsgh",
		Port: 8000,
	}

	dpl := manifestdto.Deployment{
		ApiVersion: "v1",
		Name: "test-deploy",
		Replicas: 3,
		Labels: map[string]string{"app": "nginx"},
		ContainerSpec: spec,
	}

	fmt.Println(manifestgenerator.CreateDeploymentYAML(&dpl))
}
