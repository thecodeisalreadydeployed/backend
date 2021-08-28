package workloadcontroller

import (
	"sync"
	"time"

	"github.com/thecodeisalreadydeployed/datastore"
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

func SetDeploymentState(deploymentID string, state model.DeploymentState) {
	zap.L().Sugar().Infof("Updating state of deployment ID '%s' to '%s'.", deploymentID, state)
	err := datastore.SetDeploymentState(deploymentID, state)
	if err != nil {
		zap.L().Sugar().Infof("Failed to set state of deployment ID '%s'.", deploymentID)
	}
}

func CheckDeployments() {
	for {
		deploymentsToCheck.Range(func(key, value interface{}) bool {
			deploymentID := key.(string)
			zap.L().Sugar().Infof("Checking state of deployment ID '%s'.", deploymentID)
			state := DeploymentState(deploymentID)
			if state == model.DeploymentStateBuildSucceeded {
				deploymentsToCheck.Delete(key)
				zap.L().Sugar().Infof("Deleted deployment ID '%s' from deploymentsToCheck. (Reason: %s)", deploymentID, state)
				SetDeploymentState(deploymentID, state)
			}
			return true
		})

		time.Sleep(1 * time.Minute)
	}
}
