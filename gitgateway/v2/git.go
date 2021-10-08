package gitgateway

import (
	"github.com/go-git/go-git/v5"
	"github.com/thecodeisalreadydeployed/errutil"
)

type GitGateway interface {
	Checkout(branch string) error
	OpenFile(filePath string) error
	WriteFile(filePath string, data string) error
	Commit(files []string, message string) (string, error)
	Pull() error
	Diff(oldCommit string, currentCommit string) ([]string, error)
}

type gitGateway struct {
	repo git.Repository
}

func NewGitGateway() (GitGateway, error) {
	return &gitGateway{}, nil
}

func (g *gitGateway) Checkout(branch string) error {
	return errutil.ErrNotImplemented
}

func (g *gitGateway) OpenFile(filePath string) error {
	return errutil.ErrNotImplemented
}

func (g *gitGateway) WriteFile(filePath string, data string) error {
	return errutil.ErrNotImplemented
}

func (g *gitGateway) Commit(files []string, message string) (string, error) {
	return "", errutil.ErrNotImplemented
}

func (g *gitGateway) Pull() error {
	return errutil.ErrNotImplemented
}

func (g *gitGateway) Diff(oldCommit string, currentCommit string) ([]string, error) {
	return []string{}, errutil.ErrNotImplemented
}
