package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/kanikogateway/v2"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

func NewDeployment(appID string) (*model.Deployment, error) {
	logger := zap.L().Sugar().With("appID", appID)
	_ = logger

	app, err := datastore.GetAppByID(datastore.GetDB(), appID)
	if err != nil {
		return nil, err
	}

	git, err := gitgateway.NewGitGatewayRemote(app.GitSource.RepositoryURL)
	if err != nil {
		return nil, err
	}

	commitHash, err := git.Head()
	if err != nil {
		return nil, err
	}

	gitSource := model.GitSource{
		CommitSHA:     commitHash,
		RepositoryURL: app.GitSource.RepositoryURL,
		Branch:        app.GitSource.Branch,
	}

	app.GitSource = gitSource
	app, err = datastore.SaveApp(datastore.GetDB(), app)
	if err != nil {
		return nil, err
	}

	deployment, err := datastore.SaveDeployment(datastore.GetDB(), &model.Deployment{
		AppID:              appID,
		GitSource:          app.GitSource,
		BuildConfiguration: app.BuildConfiguration,
		State:              model.DeploymentStateQueueing,
	})
	if err != nil {
		return nil, err
	}

	kaniko, err := kanikogateway.NewKanikoGateway(app.ProjectID, app.ID, deployment.ID, deployment.GitSource.RepositoryURL, deployment.GitSource.Branch, deployment.BuildConfiguration, nil)
	if err == nil {
		if err != nil {
			return nil, err
		}

		podName, err := kaniko.Deploy()
		if err != nil {
			return nil, err
		}

		_ = podName
	}

	return deployment, nil
}
