package repositoryobserver

import (
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/workloadcontroller"
	"go.uber.org/zap"
)

func hasChanges(gs *model.GitSource, gw *gitgateway.GitGateway) bool {
	old := gw.GetCommit(gs.LastObservedCommitSHA)
	current := gw.GetCurrentCommit()
	return gitgateway.HasProperDiff(old, current)
}

func ObserveGitSources() {
	apps, err := datastore.GetAllApps()
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	zap.L().Info("Observing source code...")
	var changes []model.App
	for _, app := range *apps {
		gw := gitgateway.NewGitGateway(app.GitSource.RepositoryURL)
		if hasChanges(&app.GitSource, &gw) {
			changes = append(changes, app)
		}
	}
	workloadcontroller.DeployNewRevisions(&changes)
}
