package gitapi

import (
	"net/url"

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

func NewGitAPIBackend() GitAPIBackend {
	return &gitAPIBackend{}
}

func (backend *gitAPIBackend) getGitProvider(repoURL string) gitProvider {
	u, err := url.Parse(repoURL)
	if err != nil {
		backend.logger.Error("cannot parse repository URL", zap.Error(err))
		return unknown
	}

	switch u.Hostname() {
	case string(gitHub):
		return gitHub
	case string(gitLab):
		return gitLab
	default:
		return unknown
	}
}

func (backend *gitAPIBackend) GetBranches(repoURL string) ([]string, error) {
	provider := backend.getGitProvider(repoURL)
	switch provider {
	case gitHub:
		return github.NewGitHubAPI()
	}
}
