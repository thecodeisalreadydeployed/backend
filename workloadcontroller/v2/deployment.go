package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/kanikogateway/v2"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

func NewDeployment(appID string) error {
	logger := zap.L().Sugar().With("appID", appID)
	_ = logger

	app, err := datastore.GetAppByID(datastore.GetDB(), appID)
	if err != nil {
		return err
	}

	deployment, err := datastore.SaveDeployment(datastore.GetDB(), &model.Deployment{
		AppID:              appID,
		GitSource:          app.GitSource,
		BuildConfiguration: app.BuildConfiguration,
	})
	if err != nil {
		return err
	}

	err = datastore.SetDeploymentState(datastore.GetDB(), deployment.ID, model.DeploymentStateQueueing)
	if err != nil {
		return err
	}

	kaniko, err := kanikogateway.NewKanikoGateway(app.ProjectID, app.ID, deployment.ID, deployment.GitSource.RepositoryURL, deployment.GitSource.Branch, deployment.BuildConfiguration, nil)
	if err != nil {
		return err
	}

	podName, err := kaniko.Deploy()
	if err != nil {
		return err
	}

	_ = podName

	return nil
}
