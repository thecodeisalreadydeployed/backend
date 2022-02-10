//go:generate mockgen -destination mock/controller.go . WorkloadController

package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/gitopscontroller"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

type WorkloadController interface {
	NewProject(project *model.Project) (*model.Project, error)
	NewDeployment(appID string, expectedCommitHash *string) (*model.Deployment, error)
}

type workloadController struct {
	logger           *zap.Logger
	gitOpsController gitopscontroller.GitOpsController
}

func NewWorkloadController(logger *zap.Logger, gitOpsController gitopscontroller.GitOpsController) WorkloadController {
	return &workloadController{logger: logger, gitOpsController: gitOpsController}
}
