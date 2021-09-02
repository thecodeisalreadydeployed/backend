package gitopscontroller

import (
	"errors"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/thecodeisalreadydeployed/config"
	gitgw "github.com/thecodeisalreadydeployed/gitgateway"
)

type GitOpsController struct {
	userspace *gitgw.GitGateway
	mutex     sync.Mutex
}

var controller *GitOpsController

func GetController() *GitOpsController {
	return controller
}

func Init() {
	gw := gitgw.NewGitGateway(config.DefaultUserspaceRepository)
	newGitOpsController(&gw)
}

func newGitOpsController(userspace *gitgw.GitGateway) {
	controller = &GitOpsController{
		userspace: userspace,
		mutex:     sync.Mutex{},
	}
}

func (c *GitOpsController) SetupUserspace() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := gitgw.InitRepository(config.DefaultUserspaceRepository)
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		return nil
	}
	return err
}
