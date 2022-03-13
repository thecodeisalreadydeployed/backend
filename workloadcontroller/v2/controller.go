//go:generate mockgen -destination mock/controller.go . WorkloadController

package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/clusterbackend"
	containerregistry "github.com/thecodeisalreadydeployed/containerregistry/types"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitopscontroller"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

type WorkloadController interface {
	NewProject(project *model.Project, dataStore datastore.DataStore) (*model.Project, error)
	NewApp(app *model.App, dataStore datastore.DataStore) (*model.App, error)
	NewDeployment(appID string, expectedCommitHash *string, dataStore datastore.DataStore) (*model.Deployment, error)
	ObserveWorkloads(datastore datastore.DataStore)
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
