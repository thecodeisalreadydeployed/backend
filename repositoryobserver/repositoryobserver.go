package repositoryobserver

import (
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/workloadcontroller"
	"go.uber.org/zap"
)

func hasChanges(gs *model.GitSource, gw *gitgateway.GitGateway) bool {
	old := gw.GetCommit(gs.CommitSHA)
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
	changes := make(map[string]string)
	for _, app := range *apps {
		gw := gitgateway.NewGitGateway(app.GitSource.RepositoryURL)
		if hasChanges(&app.GitSource, &gw) {
			changes[app.ID] = app.GitSource.CommitSHA
		}
	}
	workloadcontroller.OnGitSourceUpdate(&changes)
}
