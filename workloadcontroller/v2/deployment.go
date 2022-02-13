package workloadcontroller

import (
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/kanikogateway"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

func (ctrl *workloadController) NewDeployment(appID string, expectedCommitHash *string) (*model.Deployment, error) {
	app, err := datastore.GetAppByID(datastore.GetDB(), appID)
	if err != nil {
		return nil, err
	}

	var commitHash string

	if expectedCommitHash == nil {
		git, err := gitgateway.NewGitGatewayRemote(app.GitSource.RepositoryURL)
		if err != nil {
			return nil, err
		}

		err = git.Checkout(app.GitSource.Branch)
		if err != nil {
			return nil, err
		}

		commitHash, err = git.Head()
		if err != nil {
			return nil, err
		}
	} else {
		commitHash = *expectedCommitHash
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

	kaniko, err := kanikogateway.NewKanikoGateway(ctrl.logger.With(zap.String("appID", appID)), ctrl.clusterBackend, app.ProjectID, app.ID, deployment.ID, deployment.GitSource.RepositoryURL, deployment.GitSource.Branch, deployment.BuildConfiguration, nil)
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
