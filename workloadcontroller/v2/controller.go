//go:generate mockgen -destination mock/controller.go . WorkloadController

package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/gitopscontroller"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor/v2"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/util"
	"go.uber.org/zap"
)

type WorkloadController interface {
	NewProject(project *model.Project) (*model.Project, error)
	NewApp(app *model.App) (*model.App, error)
	NewDeployment(appID string, expectedCommitHash *string) (*model.Deployment, error)
	ObserveWorkloads()
}

type workloadController struct {
	logger           *zap.Logger
	gitOpsController gitopscontroller.GitOpsController
	kubernetesClient *kubernetesinteractor.KubernetesInteractor
}

func NewWorkloadController(logger *zap.Logger, gitOpsController gitopscontroller.GitOpsController) WorkloadController {
	kubernetesClient, err := kubernetesinteractor.NewKubernetesInteractor()
	if err != nil {
		if util.IsProductionEnvironment() || util.IsKubernetesTestEnvironment() {
			logger.Warn("cannot create Kubernetes client", zap.Error(err))
		} else {
			logger.Warn("cannot create Kubernetes client", zap.Error(err))
		}
	}

	return &workloadController{logger: logger, gitOpsController: gitOpsController, kubernetesClient: &kubernetesClient}
}
