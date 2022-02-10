//go:generate mockgen -destination mock/gitopscontroller.go . GitOpsController
package gitopscontroller

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/gitopscontroller/argocd"
	"github.com/thecodeisalreadydeployed/gitopscontroller/kustomize"
	"github.com/thecodeisalreadydeployed/manifestgenerator"
	"github.com/thecodeisalreadydeployed/util"
	"go.uber.org/zap"
)

type GitOpsController interface {
	SetupProject(projectID string) error
	SetupApp(projectID string, appID string) error
	SetContainerImage(projectID string, appID string, newImage string) error
}

type gitOpsController struct {
	logger       *zap.Logger
	user         gitgateway.GitGateway
	path         string
	argoCDClient argocd.ArgoCDClient
}

var once sync.Once
var mutex sync.Mutex

func setupUserspace() {
	once.Do(func() {
		path := config.DefaultUserspaceRepository()
		gateway, err := gitgateway.NewGitRepository(path)
		if err != nil {
			panic(err)
		}

		writeErr := gateway.WriteFile("kustomization.yml", "")
		if writeErr != nil {
			panic("cannot write kustomization.yml")
		}

		_, commitErr := gateway.Commit([]string{"kustomization.yml"}, "initial commit")
		if commitErr != nil {
			panic(commitErr)
		}
	})
}

func NewGitOpsController(logger *zap.Logger) GitOpsController {
	setupUserspace()
	path := config.DefaultUserspaceRepository()
	userspace, err := gitgateway.NewGitGatewayLocal(path)
	if err != nil {
		panic(err)
	}

	argoCDClient := argocd.NewArgoCDClient(logger.With(zap.String("userspace", path)), "userspace", path)
	if err = argoCDClient.CreateApp(); err != nil {
		if util.IsProductionEnvironment() || util.IsKubernetesTestEnvironment() {
			logger.Fatal("cannot create Argo CD application", zap.Error(err))
		} else {
			logger.Warn("cannot create Argo CD application", zap.Error(err))
		}
	}

	return &gitOpsController{logger: logger, user: userspace, path: path, argoCDClient: argoCDClient}
}

func (g *gitOpsController) SetupProject(projectID string) error {
	mutex.Lock()
	defer mutex.Unlock()

	kustomizationFile := fmt.Sprintf("%s/kustomization.yml", projectID)
	writeErr := g.user.WriteFile(kustomizationFile, "")
	if writeErr != nil {
		return errors.New("cannot write kustomization.yml")
	}

	kustomizeErr := kustomize.AddResources(filepath.Join(g.path, "kustomization.yml"), []string{kustomizationFile})
	if kustomizeErr != nil {
		return errors.New("cannot write kustomization.yml")
	}

	commitHash, commitErr := g.user.Commit([]string{"kustomization.yml", kustomizationFile}, projectID)
	if commitErr != nil {
		return commitErr
	}

	fmt.Printf("commitHash: %v\n", commitHash)

	err := g.argoCDClient.Sync()
	if err != nil {
		return err
	}

	return nil
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
			"beta.deploys.dev/project-id": projectID,
			"beta.deploys.dev/app-id":     appID,
			"beta.deploys.dev/part-of":    "gitopscontroller",
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
			"beta.deploys.dev/project-id": projectID,
			"beta.deploys.dev/app-id":     appID,
			"beta.deploys.dev/part-of":    "gitopscontroller",
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

	kustomizeErr := kustomize.AddResources(filepath.Join(g.path, kustomizationFile), []string{"deployment.yml", "service.yml"})
	if kustomizeErr != nil {
		return errors.New("cannot write kustomization.yml")
	}

	kustomizeErr = kustomize.AddResources(filepath.Join(g.path, projectID, "kustomization.yml"), []string{filepath.Join(appID, "kustomization.yml")})
	if kustomizeErr != nil {
		return errors.New("cannot write kustomization.yml")
	}

	commit, commitErr := g.user.Commit([]string{kustomizationFile, deploymentFile, serviceFile, filepath.Join(projectID, "kustomization.yml")}, prefix)
	if commitErr != nil {
		return commitErr
	}

	fmt.Printf("commit: %v\n", commit)

	err := g.argoCDClient.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (g *gitOpsController) SetContainerImage(projectID string, appID string, newImage string) error {
	mutex.Lock()
	defer mutex.Unlock()

	prefix := fmt.Sprintf("%s/%s", projectID, appID)
	kustomizationFile := filepath.Join(prefix, "kustomization.yml")
	err := kustomize.SetImage(filepath.Join(g.path, kustomizationFile), "codedeploy://"+appID, newImage)
	if err != nil {
		return err
	}

	_, commitErr := g.user.Commit([]string{kustomizationFile}, fmt.Sprintf("%s: %s", prefix, newImage))
	if commitErr != nil {
		return commitErr
	}
	return g.argoCDClient.Sync()
}
