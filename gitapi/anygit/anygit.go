package anygit

import (
	"github.com/thecodeisalreadydeployed/gitapi/provider"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
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
		return []string{}, gitErr
	}
	return git.GetBranches()
}

func (api *anyGitAPI) GetFiles(branch string) ([]string, error) {
	git, gitErr := gitgateway.NewGitGatewayRemote(api.repoURL)
	if gitErr != nil {
		return []string{}, gitErr
	}
	return git.GetFiles(branch)
}

func (api *anyGitAPI) GetRaw(branch string, path string) (string, error) {
	git, gitErr := gitgateway.NewGitGatewayRemote(api.repoURL)
	if gitErr != nil {
		return "", gitErr
	}
	return git.GetRaw(branch, path)
}

func (api *anyGitAPI) FillGitSource(gs model.GitSource) (model.GitSource, error) {
	gs, err := gitgateway.Info(gs.RepositoryURL, gs.Branch)
	if err != nil {
		return model.GitSource{}, err
	}
	return gs, nil
}
