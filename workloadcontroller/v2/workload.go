package workloadcontroller

import (
	"time"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

func (ctrl *workloadController) ObserveWorkloads() {
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
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}
