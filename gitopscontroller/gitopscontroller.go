package gitopscontroller

import (
	"errors"
	"fmt"
	"sync"

	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/manifestgenerator"
)

type GitOpsController interface {
	SetupProject(projectID string) error
	SetupApp(projectID string, appID string) error
	UpdateContainerImage(projectID string, appID string, newImage string) error
}

type gitOpsController struct {
	user gitgateway.GitGateway
}

var once sync.Once
var mutex sync.Mutex

func setupUserspace() {
	once.Do(func() {
		path := config.DefaultUserspaceRepository
		gateway, err := gitgateway.NewGitRepository(path)
		if err != nil {
			panic(err)
		}
		if err = gateway.Checkout("main"); err != nil {
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

	return &gitOpsController{user: userspace}
}

func (g *gitOpsController) SetupProject(projectID string) error {
	return errutil.ErrNotImplemented
}

func (g *gitOpsController) SetupApp(projectID string, appID string) error {
	mutex.Lock()
	defer mutex.Unlock()

	prefix := fmt.Sprintf("%s/%s", projectID, appID)

	kustomizationFile := fmt.Sprintf("%s/kustomization.yml", prefix)
	deploymentFile := fmt.Sprintf("%s/deployment.yml", prefix)
	serviceFile := fmt.Sprintf("%s/service.yml", prefix)

	writeErr := g.user.WriteFile(kustomizationFile, "")
	if writeErr != nil {
		return errors.New("cannot write kustomization.yml")
	}

	deploymentYAML, generateErr := manifestgenerator.GenerateDeploymentYAML(&manifestgenerator.GenerateDeploymentOptions{
		Name:      appID,
		Namespace: projectID,
		Labels: map[string]string{
			"project.api.deploys.dev/id": projectID,
			"app.api.deploys.dev/id":     appID,
			"api.deploys.dev/part-of":    "gitopscontroller",
		},
		ContainerImage: "codedeploy://" + appID,
	})

	if generateErr != nil {
		return errors.New("cannot generate deployment.yml")
	}

	serviceYAML, generateErr := manifestgenerator.GenerateServiceYAML(&manifestgenerator.GenerateServiceOptions{
		Name:      appID,
		Namespace: projectID,
		Labels: map[string]string{
			"project.api.deploys.dev/id": projectID,
			"app.api.deploys.dev/id":     appID,
			"api.deploys.dev/part-of":    "gitopscontroller",
		},
	})

	if generateErr != nil {
		return errors.New("cannot generate service.yml")
	}

	writeErr = g.user.WriteFile(deploymentFile, deploymentYAML)
	if writeErr != nil {
		return errors.New("cannot write deployment.yml")
	}

	writeErr = g.user.WriteFile(serviceFile, serviceYAML)
	if writeErr != nil {
		return errors.New("cannot write service.yml")
	}

	commit, commitErr := g.user.Commit([]string{kustomizationFile, deploymentFile, serviceFile}, prefix)
	if commitErr != nil {
		return commitErr
	}

	fmt.Printf("commit: %v\n", commit)

	return nil
}

func (g *gitOpsController) UpdateContainerImage(projectID string, appID string, newImage string) error {
	return errutil.ErrNotImplemented
}
