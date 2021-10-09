package repositoryobserver

import (
	"fmt"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
)

func ObserveGitSources() {
	apps, err := datastore.GetObservableApps()
	if err != nil {
		return
	}

	for _, app := range *apps {
		go func(gitSource model.GitSource) {
			check(gitSource.RepositoryURL, gitSource.Branch, gitSource.CommitSHA)
		}(app.GitSource)
	}
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
