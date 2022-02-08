package anygit

import (
	"github.com/thecodeisalreadydeployed/gitapi/provider"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"go.uber.org/zap"
)

type anyGitAPI struct {
	logger  *zap.Logger
	repoURL string
}

func NewAnyGitAPI(logger *zap.Logger, repoURL string) provider.GitProvider {
	return &anyGitAPI{logger: logger, repoURL: repoURL}
}

func (api *anyGitAPI) GetBranches() ([]string, error) {
	git, gitErr := gitgateway.NewGitGatewayRemote(api.repoURL)
	if gitErr != nil {
		return []string{}, nil
	}
	return git.GetBranches()
}

func (api *anyGitAPI) GetFiles(branch string) ([]string, error) {
	return []string{}, nil
}

func (api *anyGitAPI) GetRaw(branch string, path string) (string, error) {
	return "", nil
}
