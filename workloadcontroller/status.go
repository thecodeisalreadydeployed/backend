package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/kanikointeractor"
	"github.com/thecodeisalreadydeployed/model"
)

func DeploymentState(deploymentID string) model.DeploymentState {
	return kanikointeractor.DeploymentState(deploymentID)
}
