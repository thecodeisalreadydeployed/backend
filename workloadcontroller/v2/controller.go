//go:generate mockgen -destination mock/controller.go . WorkloadController

package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/clusterbackend"
	containerregistry "github.com/thecodeisalreadydeployed/containerregistry/types"
	"github.com/thecodeisalreadydeployed/gitopscontroller"
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
	logger            *zap.Logger
	gitOpsController  gitopscontroller.GitOpsController
	clusterBackend    clusterbackend.ClusterBackend
	containerRegistry containerregistry.ContainerRegistry
}

func NewWorkloadController(
	logger *zap.Logger,
	gitOpsController gitopscontroller.GitOpsController,
	clusterBackend clusterbackend.ClusterBackend,
	containerRegistry containerregistry.ContainerRegistry,
) WorkloadController {
	return &workloadController{
		logger:            logger,
		gitOpsController:  gitOpsController,
		clusterBackend:    clusterBackend,
		containerRegistry: containerRegistry,
	}
}
