package gitgateway

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/errutil"
	"go.uber.org/zap"
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
		fmt.Printf("worktreeErr: %v\n", worktreeErr)
		zap.L().Error(worktreeErr.Error())
		return errutil.ErrFailedPrecondition
	}

	_, statErr := w.Filesystem.Stat(filePath)
	if os.IsNotExist(statErr) {
		_, createErr := w.Filesystem.Create(filePath)
		if createErr != nil {
			fmt.Printf("createErr: %v\n", createErr)
			return errutil.ErrFailedPrecondition
		}
	}

	f, openErr := w.Filesystem.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, defaultMode)
	if openErr != nil {
		fmt.Printf("openErr: %v\n", openErr)
		zap.L().Error(openErr.Error())
		return errutil.ErrFailedPrecondition
	}

	_, writeErr := f.Write([]byte(data))
	if writeErr != nil {
		fmt.Printf("writeErr: %v\n", writeErr)
		zap.L().Error(writeErr.Error())
		return errutil.ErrFailedPrecondition
	}

	return nil
}

func (g *gitGateway) Commit(files []string, message string) (string, error) {
	w, worktreeErr := g.repo.Worktree()
	if worktreeErr != nil {
		return "", errutil.ErrFailedPrecondition
	}

	for _, f := range files {
		_, addErr := w.Add(f)
		if addErr != nil {
			return "", errutil.ErrFailedPrecondition
		}
	}

	commit, commitErr := w.Commit(message, &git.CommitOptions{
		Author: config.DefaultGitSignature(),
	})

	if commitErr != nil {
		return "", errutil.ErrFailedPrecondition
	}

	commitHash := commit.String()
	zap.L().Info(commitHash)

	return commitHash, nil
}

func (g *gitGateway) Pull() error {
	return errutil.ErrNotImplemented
}

func (g *gitGateway) Diff(oldCommit string, currentCommit string) ([]string, error) {
	return []string{}, errutil.ErrNotImplemented
}
