package workloadcontroller

import (
	"sync"
	"time"

	"github.com/thecodeisalreadydeployed/kanikointeractor"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

var deploymentsToCheck = sync.Map{} // [deploymentID]struct{}

func DeploymentState(deploymentID string) model.DeploymentState {
	return kanikointeractor.DeploymentState(deploymentID)
}

func AddDeploymentToCheck(deploymentID string) {
	deploymentsToCheck.Store(deploymentID, struct{}{})
	zap.L().Sugar().Infof("Added deployment ID '%s' to deploymentsToCheck.", deploymentID)
}

func CheckDeployments() {
	for {
		deploymentsToCheck.Range(func(key, value interface{}) bool {
			deploymentID := key.(string)
			zap.L().Sugar().Infof("Checking status of deployment ID '%s'.", deploymentID)
			state := DeploymentState(deploymentID)
			if state == model.DeploymentStateBuildSucceeded {
				deploymentsToCheck.Delete(key)
			}
			return true
		})

		time.Sleep(1 * time.Minute)
	}
}
