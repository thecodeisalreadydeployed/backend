package workloadcontroller

import (
	"time"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/util"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
)

func (ctrl *workloadController) ObserveWorkloads() {
	if util.IsDevEnvironment() || util.IsDockerTestEnvironment() {
		return
	}

	for {
		pendingDeployments, err := datastore.GetPendingDeployments(datastore.GetDB())
		if err != nil {
			ctrl.logger.Error(err.Error())
			continue
		}

		for _, deployment := range *pendingDeployments {
			if deployment.State == model.DeploymentStateBuilding {
				pods, err := ctrl.clusterBackend.Pods("codedeploy-internal", map[string]string{
					"beta.deploys.dev/deployment-id": deployment.ID,
					"beta.deploys.dev/component":     "imagebuilder",
				})

				if err != nil {
					ctrl.logger.Error(err.Error())
					continue
				}

				for _, p := range pods {
					ctrl.logger.Debug(p.Name, zap.String("phase", string(p.Status.Phase)), zap.String("selfLink", p.SelfLink), zap.String("startTime", p.Status.StartTime.String()))
					switch p.Status.Phase {
					case v1.PodSucceeded:
						datastore.SetDeploymentState(datastore.GetDB(), deployment.ID, model.DeploymentStateBuildSucceeded)
					case v1.PodFailed:
						datastore.SetDeploymentState(datastore.GetDB(), deployment.ID, model.DeploymentStateError)
					}
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}
