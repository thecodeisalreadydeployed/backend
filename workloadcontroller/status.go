package workloadcontroller

import (
	"sync"

	"github.com/thecodeisalreadydeployed/kanikointeractor"
	"github.com/thecodeisalreadydeployed/model"
)

var deploymentsToCheck = sync.Map{} // [deploymentID]struct{}

func DeploymentState(deploymentID string) model.DeploymentState {
	return kanikointeractor.DeploymentState(deploymentID)
}

func CheckDeployments() {
	for {
		deploymentsToCheck.Range(func(key, value interface{}) bool {
			deploymentID := key.(string)
			state := DeploymentState(deploymentID)
			if state == model.DeploymentStateBuildSucceeded {
				deploymentsToCheck.Delete(key)
			}
			return true
		})
	}
}
