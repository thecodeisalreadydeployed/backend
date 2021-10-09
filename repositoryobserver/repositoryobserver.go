package repositoryobserver

import (
	"fmt"
	"sync"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
)

func ObserveGitSources() {
	apps, err := datastore.GetObservableApps()
	if err != nil {
		return
	}

	_ = checkGitSources(*apps)
}

func checkGitSources(apps []model.App) map[string]string {
	_appsToUpdate := sync.Map{}

	for _, app := range apps {
		go func(appID string, gitSource model.GitSource) {
			commit := check(gitSource.RepositoryURL, gitSource.Branch, gitSource.CommitSHA)
			if commit != nil {
				_appsToUpdate.Store(appID, commit)
			}
		}(app.ID, app.GitSource)
	}
	appsToUpdate := map[string]string{}
	_appsToUpdate.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		appsToUpdate[k] = v
		return true
	})

	return appsToUpdate
}

func check(repoURL string, branch string, currentCommitSHA string) *string {
	git, err := gitgateway.NewGitGatewayRemote(repoURL)
	if err != nil {
		return nil
	}

	checkoutErr := git.Checkout(branch)
	if checkoutErr != nil {
		return nil
	}

	ref, err := git.Head()
	if err != nil {
		return nil
	}

	diff, diffErr := git.Diff(currentCommitSHA, ref)
	if diffErr != nil {
		return nil
	}

	if len(diff) > 0 {
		fmt.Printf("len(diff): %v\n", len(diff))
		return &ref
	}

	return nil
}
