package gitopscontroller

import (
	"fmt"
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
	u     gitgateway.GitGateway
	mutex sync.Mutex
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

	return &gitOpsController{u: userspace, mutex: sync.Mutex{}}
}

func (g *gitOpsController) SetupProject(projectID string) error {
	return errutil.ErrNotImplemented
}

func (g *gitOpsController) SetupApp(projectID string, appID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	prefix := fmt.Sprintf("%s/%s", projectID, appID)

	kustomizationFile := fmt.Sprintf("%s/kustomization.yml", prefix)
	deploymentFile := fmt.Sprintf("%s/deployment.yml", prefix)
	serviceFile := fmt.Sprintf("%s/service.yml", prefix)

	writeErr := g.u.WriteFile(kustomizationFile, "")
	if writeErr != nil {
		return errutil.ErrFailedPrecondition
	}

	writeErr = g.u.WriteFile(deploymentFile, "")
	if writeErr != nil {
		return errutil.ErrFailedPrecondition
	}

	writeErr = g.u.WriteFile(serviceFile, "")
	if writeErr != nil {
		return errutil.ErrFailedPrecondition
	}

	commit, commitErr := g.u.Commit([]string{kustomizationFile, deploymentFile, serviceFile}, prefix)
	if commitErr != nil {
		return errutil.ErrFailedPrecondition
	}

	fmt.Printf("commit: %v\n", commit)

	return nil
}

func (g *gitOpsController) UpdateContainerImage(projectID string, appID string, newImage string) error {
	return errutil.ErrNotImplemented
}
