package gitopscontroller

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/gitinteractor"
)

type GitOpsController struct {
	Userspace *gitinteractor.GitInteractor
}

func (c *GitOpsController) Init() {
	gi := gitinteractor.NewGitInteractor(config.DefaultUserspaceRepository)
	c.Userspace = &gi
}

func SetupUserspace() error {
	err := gitinteractor.InitRepository(config.DefaultUserspaceRepository)
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		return nil
	}
	return err
}
