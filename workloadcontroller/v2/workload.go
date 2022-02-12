package workloadcontroller

import (
	"fmt"
	"time"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
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
					fmt.Printf("p.Status.Phase: %v\n", p.Status.Phase)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}
