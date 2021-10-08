package gitgateway

import (
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/errutil"
)

type GitGateway interface {
	Checkout(branch string) error
	OpenFile(filePath string) (string, error)
	WriteFile(filePath string, data string) error
	Commit(files []string, message string) (string, error)
	Pull() error
	Diff(oldCommit string, currentCommit string) ([]string, error)
}

type gitGateway struct {
	repo *git.Repository
}

func NewGitGatewayLocal(path string) (GitGateway, error) {
	repo, openErr := git.PlainOpen(path)

	if openErr != nil {
		return nil, errutil.ErrFailedPrecondition
	}

	return &gitGateway{repo: repo}, nil
}

func NewGitGatewayRemote(url string) (GitGateway, error) {
	repo, cloneErr := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})

	if cloneErr != nil {
		return nil, errutil.ErrFailedPrecondition
	}

	return &gitGateway{repo: repo}, nil
}

func (g *gitGateway) Checkout(branch string) error {
	w, worktreeErr := g.repo.Worktree()
	if worktreeErr != nil {
		return errutil.ErrFailedPrecondition
	}

	checkoutErr := w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})

	if checkoutErr != nil {
		return errutil.ErrFailedPrecondition
	}

	return nil
}

func (g *gitGateway) OpenFile(filePath string) (string, error) {
	defaultMode := os.FileMode(0666)

	w, worktreeErr := g.repo.Worktree()
	if worktreeErr != nil {
		return "", errutil.ErrFailedPrecondition
	}

	f, openErr := w.Filesystem.OpenFile(filePath, os.O_RDONLY, defaultMode)
	if openErr != nil {
		return "", errutil.ErrFailedPrecondition
	}

	read, readErr := ioutil.ReadAll(f)
	if readErr != nil {
		return "", errutil.ErrFailedPrecondition
	}

	return cast.ToString(read), nil
}

func (g *gitGateway) WriteFile(filePath string, data string) error {
	defaultMode := os.FileMode(0666)

	w, worktreeErr := g.repo.Worktree()
	if worktreeErr != nil {
		return errutil.ErrFailedPrecondition
	}

	f, openErr := w.Filesystem.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, defaultMode)
	if openErr != nil {
		return errutil.ErrFailedPrecondition
	}

	_, writeErr := f.Write([]byte(data))
	if writeErr != nil {
		return errutil.ErrFailedPrecondition
	}

	return nil
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
