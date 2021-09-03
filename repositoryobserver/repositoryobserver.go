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

func ObserveGitSource() {
	zap.L().Info("Observing source code...")

	_, err := datastore.GetAllApps()
	if err != nil {
		zap.L().Error("Error while observing apps. Cannot access database.")
		return
	}

	sc := gitgateway.NewGitGateway("/home/jean/Desktop/gittest")

	gs := model.GitSource{
		Provider:              "",
		Organization:          "",
		CommitSHA:             "",
		CommitMessage:         "",
		CommitAuthorName:      "",
		RepositoryURL:         "",
		LastObservedCommitSHA: "a1d95e5b2ac18b8e1ad713d39de6e57a2479d4e2",
	}

	if hasChanges(&gs, &sc) {
		zap.L().Info("Source code has changed. Deploying new revision...")
		workloadcontroller.OnGitSourceUpdate(true)
	}

	// TODO: App's source code should be valid.
	//for _, app := range *apps {
	//	if hasChanges(&app.GitSource) {
	//		zap.L().Info("Source code has changed. Deploying new revision...")
	//		workloadcontroller.DeployNewRevision()
	//	}
	//}
}
