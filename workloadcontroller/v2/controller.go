//go:generate mockgen -destination mock/controller.go . WorkloadController

package workloadcontroller

import "github.com/thecodeisalreadydeployed/model"

type WorkloadController interface {
	NewDeployment(appID string, expectedCommitHash *string) (*model.Deployment, error)
}

type workloadController struct{}

func NewWorkloadController() WorkloadController {
	return &workloadController{}
}
