package gitapi

import (
	"fmt"
	"net/url"
	"strings"

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
