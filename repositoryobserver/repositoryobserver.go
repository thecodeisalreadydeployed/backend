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
	appsToUpdate := sync.Map{}
	var wg sync.WaitGroup

	expected := map[string]string{}
	for _, app := range apps {
		wg.Add(1)
		go func(appID string, gitSource model.GitSource) {
			defer wg.Done()
			commit := check(gitSource.RepositoryURL, gitSource.Branch, gitSource.CommitSHA)
			if commit != nil {
				fmt.Printf("appID: %v\n", appID)
				fmt.Printf("commit: %v\n", *commit)
				appsToUpdate.Store(appID, *commit)
			}
		}(app.ID, app.GitSource)
	}

	wg.Wait()
	appsToUpdate.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		expected[k] = v
		return true
	})

	return expected
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
