package gitopscontroller

import (
	"errors"
	"path/filepath"
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

func (c *GitOpsController) Write(path string, data string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	dst := filepath.Join("/", path)

	file := filepath.Base(dst)
	dir := filepath.Dir(dst)
	c.userspace.WriteFile(dir, file, []byte(data))
	c.userspace.Add(dst)
	c.userspace.Commit(dst)
}
