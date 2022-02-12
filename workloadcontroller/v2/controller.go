//go:generate mockgen -destination mock/controller.go . WorkloadController

package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/gitopscontroller"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor/v2"
	"github.com/thecodeisalreadydeployed/model"
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
	return &workloadController{logger: logger, gitOpsController: gitOpsController, kubernetesClient: nil}
}
