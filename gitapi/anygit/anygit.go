package anygit

import (
	"github.com/thecodeisalreadydeployed/gitapi/provider"
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
	return []string{}, nil
}

func (api *anyGitAPI) GetFiles(branch string) ([]string, error) {
	return []string{}, nil
}

func (api *anyGitAPI) GetRaw(branch string, path string) (string, error) {
	return "", nil
}
