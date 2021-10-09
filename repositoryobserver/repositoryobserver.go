package repositoryobserver

import (
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
)

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
		return &ref
	}

	return nil
}
