package gitopscontroller

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/thecodeisalreadydeployed/config"
	gitgw "github.com/thecodeisalreadydeployed/gitgateway"
)

type GitOpsController struct {
	userspace *gitgw.GitGateway
}

func (c *GitOpsController) Init() {
	gw := gitgw.NewGitGateway(config.DefaultUserspaceRepository)
	c.userspace = &gw
}

func SetupUserspace() error {
	err := gitgw.InitRepository(config.DefaultUserspaceRepository)
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		return nil
	}
	return err
}
