package workloadcontroller

import (
	"fmt"
	"strings"
	"time"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/util"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
)

func (ctrl *workloadController) setContainerImage(appID string, deploymentID string, dataStore datastore.DataStore) {
	app, err := dataStore.GetAppByID(appID)
	if err != nil {
		ctrl.logger.Error(err.Error(), zap.String("appID", appID), zap.String("deploymentID", deploymentID))
		return
	}

	newImage := ctrl.containerRegistry.RegistryFormat(app.ID, deploymentID)
	if util.IsKubernetesTestEnvironment() {
		newImage = fmt.Sprintf("localhost:31500/%s:%s", app.ID, deploymentID)
	}

	ctrl.logger.Info("setting container image", zap.String("appID", appID), zap.String("deploymentID", deploymentID), zap.String("newImage", newImage))
	err = ctrl.gitOpsController.SetContainerImage(app.ProjectID, app.ID, deploymentID, newImage)
	if err != nil {
		ctrl.logger.Error("cannot set container image", zap.Error(err), zap.String("appID", appID), zap.String("deploymentID", deploymentID))
		return
	}
}

func (ctrl *workloadController) ObserveWorkloads(dataStore datastore.DataStore) {
	if util.IsDevEnvironment() || util.IsDockerTestEnvironment() {
		return
	}

	for {
		pendingDeployments, err := dataStore.GetPendingDeployments()
		if err != nil {
			ctrl.logger.Error(err.Error())
			continue
		}

		numberOfPendingDeployments := len(*pendingDeployments)
		if numberOfPendingDeployments != 0 {
			ctrl.logger.Debug("number of pending deployments is " + fmt.Sprint(numberOfPendingDeployments))
		}

		for _, deployment := range *pendingDeployments {
			ctrl.logger.Debug("processing deployment", zap.Any("deployment", deployment))

			if deployment.State == model.DeploymentStateQueueing {
				timeLimit := deployment.UpdatedAt.Add(15 * time.Minute)
				if time.Now().After(timeLimit) {
					err = dataStore.SetDeploymentState(deployment.ID, model.DeploymentStateError)
					if err != nil {
						ctrl.logger.Error(
							"cannot set deployment state",
							zap.String("deploymentID", deployment.ID),
							zap.String("desiredState", string(model.DeploymentStateError)),
							zap.Error(err),
						)
					}
				}
			}

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
						ctrl.logger.Info(p.Name, zap.Any("pod", p), zap.Any("containerStatuses", p.Status.ContainerStatuses))
						err = dataStore.SetDeploymentState(deployment.ID, model.DeploymentStateBuildSucceeded)
						if err != nil {
							ctrl.logger.Error(
								"cannot set deployment state",
								zap.String("deploymentID", deployment.ID),
								zap.String("desiredState", string(model.DeploymentStateBuildSucceeded)),
								zap.String("podSelfLink", p.SelfLink),
								zap.Error(err),
							)
						}
						err = ctrl.clusterBackend.DeletePod(p.Namespace, p.Name)
						if err != nil {
							ctrl.logger.Error(
								"cannot delete Pod",
								zap.String("deploymentID", deployment.ID),
								zap.String("podSelfLink", p.SelfLink),
								zap.Error(err),
							)
						}
						ctrl.setContainerImage(deployment.AppID, deployment.ID, dataStore)
						err = dataStore.SetDeploymentState(deployment.ID, model.DeploymentStateCommitted)
						if err != nil {
							ctrl.logger.Error(
								"cannot set deployment state",
								zap.String("deploymentID", deployment.ID),
								zap.String("desiredState", string(model.DeploymentStateCommitted)),
								zap.String("podSelfLink", p.SelfLink),
								zap.Error(err),
							)
						}
					case v1.PodFailed:
						err = dataStore.SetDeploymentState(deployment.ID, model.DeploymentStateError)
						if err != nil {
							ctrl.logger.Error(
								"cannot set deployment state",
								zap.String("deploymentID", deployment.ID),
								zap.String("desiredState", string(model.DeploymentStateError)),
								zap.String("podSelfLink", p.SelfLink),
								zap.Error(err),
							)
						}
					}
				}
			}

			if deployment.State == model.DeploymentStateCommitted {
				// activeDeploymentID, err := ctrl.statusAPIBackend.GetActiveDeploymentID(deployment.AppID, dataStore)
				// if err != nil {
				// 	ctrl.logger.Error(err.Error())
				// 	continue
				// }

				// if activeDeploymentID != deployment.ID {
				// 	continue
				// }

				pods, err := ctrl.clusterBackend.Pods("", map[string]string{
					"beta.deploys.dev/app-id": deployment.AppID,
				})

				if err != nil {
					ctrl.logger.Error(err.Error())
					continue
				}

				for _, p := range pods {
					ctrl.logger.Debug(p.Name, zap.String("phase", string(p.Status.Phase)), zap.String("selfLink", p.SelfLink), zap.String("startTime", p.Status.StartTime.String()))

					if p.Status.Phase == v1.PodRunning {
						err = dataStore.SetDeploymentState(deployment.ID, model.DeploymentStateReady)
						if err != nil {
							ctrl.logger.Error(
								"cannot set deployment state",
								zap.String("deploymentID", deployment.ID),
								zap.String("desiredState", string(model.DeploymentStateError)),
								zap.String("podSelfLink", p.SelfLink),
								zap.Error(err),
							)
						}
						break
					}

					if p.Status.Phase == v1.PodFailed {
						for _, container := range p.Spec.Containers {
							if container.Name == "container0" && !strings.HasPrefix(container.Image, "codedeploy://") {
								err = dataStore.SetDeploymentState(deployment.ID, model.DeploymentStateError)
								if err != nil {
									ctrl.logger.Error(
										"cannot set deployment state",
										zap.String("deploymentID", deployment.ID),
										zap.String("desiredState", string(model.DeploymentStateError)),
										zap.String("podSelfLink", p.SelfLink),
										zap.Error(err),
									)
								}
							}
						}
					}
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}
