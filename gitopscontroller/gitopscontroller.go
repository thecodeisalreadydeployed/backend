package gitopscontroller

import (
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/gitinteractor"
)

type GitOpsController struct {
	gitInteractor *gitinteractor.GitInteractor
}

func (c *GitOpsController) Init() {
	git := gitinteractor.NewGitInteractor()
	c.gitInteractor = &git
}

func SetupUserspace() error {
	err := gitinteractor.InitRepository(config.DefaultUserspaceRepository)
	return err
}
