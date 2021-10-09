package gitopscontroller

import (
	"sync"

	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
)

type GitOpsController interface {
	SetupProject(projectID string) error
	SetupApp(projectID string, appID string) error
	UpdateContainerImage(projectID string, appID string, newImage string) error
}

type gitOpsController struct {
	userspace *gitgateway.GitGateway
	mutex     sync.Mutex
}

var once sync.Once

func setupUserspace() {
	once.Do(func() {
		path := config.DefaultUserspaceRepository
		_, err := gitgateway.NewGitRepository(path)
		if err != nil {
			panic(err)
		}
	})
}

func NewGitOpsController() GitOpsController {
	setupUserspace()

	userspace, err := gitgateway.NewGitGatewayLocal(config.DefaultUserspaceRepository)
	if err != nil {
		panic(err)
	}

	return &gitOpsController{userspace: &userspace, mutex: sync.Mutex{}}
}

func (g *gitOpsController) SetupProject(projectID string) error {
	return errutil.ErrNotImplemented
}

func (g *gitOpsController) SetupApp(projectID string, appID string) error {
	return errutil.ErrNotImplemented
}

func (g *gitOpsController) UpdateContainerImage(projectID string, appID string, newImage string) error {
	return errutil.ErrNotImplemented
}
