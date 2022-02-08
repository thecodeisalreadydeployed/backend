package gitapi

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/thecodeisalreadydeployed/gitapi/anygit"
	"github.com/thecodeisalreadydeployed/gitapi/github"
	"github.com/thecodeisalreadydeployed/gitapi/provider"
	"go.uber.org/zap"
)

type GitAPIBackend interface {
	GetBranches(repoURL string) ([]string, error)
	GetFiles(repoURL string, branch string) ([]string, error)
	GetRaw(repoURL string, branch string, path string) (string, error)
}

type gitAPIBackend struct {
	logger *zap.Logger
}

type gitProvider string

const (
	gitHub  gitProvider = "github.com"
	gitLab  gitProvider = "gitlab.com"
	unknown gitProvider = "deploys.dev"
)

func NewGitAPIBackend(logger *zap.Logger) GitAPIBackend {
	return &gitAPIBackend{logger: logger}
}

func (backend *gitAPIBackend) parseRepositoryURL(repoURL string) (*url.URL, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		backend.logger.Error("cannot parse repository URL", zap.Error(err))
		return nil, fmt.Errorf("cannot parse repository URL")
	}
	return u, nil
}

func (backend *gitAPIBackend) getOwnerAndRepo(repoURL string) (string, string) {
	u, err := backend.parseRepositoryURL(repoURL)
	if err != nil {
		return "", ""
	}
	parts := strings.Split(u.Path, "/")
	if len(parts) < 2 {
		backend.logger.Error("invalid repository URL")
		return "", ""
	}
	return parts[len(parts)-2], parts[len(parts)-1]
}

func (backend *gitAPIBackend) getGitProvider(repoURL string) gitProvider {
	u, err := backend.parseRepositoryURL(repoURL)
	if err != nil {
		return unknown
	}

	switch u.Hostname() {
	case string(gitHub):
		return gitHub
	case string(gitLab):
		return unknown
	default:
		return unknown
	}
}

func (backend *gitAPIBackend) getGitProviderAPI(repoURL string) (*provider.GitProvider, error) {
	logger := backend.logger.With(zap.String("repoURL", repoURL))
	provider := backend.getGitProvider(repoURL)
	owner, repo := backend.getOwnerAndRepo(repoURL)
	if provider != unknown {
		if len(owner) == 0 {
			logger.Error("repository owner cannot be empty")
			return nil, fmt.Errorf("repository owner cannot be empty")
		}
		if len(repo) == 0 {
			logger.Error("repository name cannot be empty")
			return nil, fmt.Errorf("repository name cannot be empty")
		}
	}
	switch provider {
	case gitHub:
		api := github.NewGitHubAPI(logger, owner, repo)
		return &api, nil
	default:
		api := anygit.NewAnyGitAPI(logger, repoURL)
		return &api, nil
	}
}

func (backend *gitAPIBackend) GetBranches(repoURL string) ([]string, error) {
	providerAPI, err := backend.getGitProviderAPI(repoURL)
	if err != nil {
		return []string{}, err
	}
	return (*providerAPI).GetBranches()
}

func (backend *gitAPIBackend) GetFiles(repoURL string, branch string) ([]string, error) {
	providerAPI, err := backend.getGitProviderAPI(repoURL)
	if err != nil {
		return []string{}, err
	}
	return (*providerAPI).GetFiles(branch)
}

func (backend *gitAPIBackend) GetRaw(repoURL string, branch string, path string) (string, error) {
	providerAPI, err := backend.getGitProviderAPI(repoURL)
	if err != nil {
		return "", err
	}
	return (*providerAPI).GetRaw(branch, path)
}
