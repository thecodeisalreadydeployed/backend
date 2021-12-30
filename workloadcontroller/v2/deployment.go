package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/datastore"
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
		AppID:     appID,
		GitSource: app.GitSource,
	})
	if err != nil {
		return err
	}

	err = datastore.SetDeploymentState(datastore.GetDB(), deployment.ID, model.DeploymentStateQueueing)
	if err != nil {
		return err
	}

	_ = deployment

	return nil
}
