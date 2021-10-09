package gitopscontroller

import (
	"sync"

	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
)

type GitOpsController interface {
	SetupUserspace()
	SetupProject(projectID string) error
	SetupApp(projectID string, appID string) error
	UpdateContainerImage(appID string, deploymentID string) error
}

type gitOpsController interface{}

var once sync.Once

func SetupUserspace() {
	once.Do(func() {
		path := config.DefaultUserspaceRepository
		_, err := gitgateway.NewGitRepository(path)
		if err != nil {
			panic(err)
		}
	})
}
