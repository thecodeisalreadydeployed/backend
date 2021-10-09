package gitgateway

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-git/go-git/v5/utils/merkletrie"
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
	Log() error
	Head() (string, error)
	Diff(oldCommit string, currentCommit string) ([]string, error)
}

type gitGateway struct {
	repo *git.Repository
}

func NewGitRepository(path string) (GitGateway, error) {
	repo, initErr := git.PlainInit(path, false)
	if initErr != nil {
		return nil, errutil.ErrFailedPrecondition
	}
	return &gitGateway{repo: repo}, nil
}

func NewGitGatewayLocal(path string) (GitGateway, error) {
	repo, openErr := git.PlainOpen(path)

	if openErr != nil {
		return nil, errutil.ErrFailedPrecondition
	}

	return &gitGateway{repo: repo}, nil
}

func NewGitGatewayRemote(url string) (GitGateway, error) {
	repo, cloneErr := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
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
		fmt.Printf("worktreeErr: %v\n", worktreeErr)
		return errutil.ErrFailedPrecondition
	}

	checkoutErr := w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
	})

	if checkoutErr != nil {
		if !errors.Is(checkoutErr, plumbing.ErrReferenceNotFound) {
			fmt.Printf("checkoutErr: %v\n", checkoutErr)
			return errutil.ErrFailedPrecondition
		} else {
			err := w.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(branch),
				Create: true,
				Force:  true,
			})

			if err != nil {
				fmt.Printf("err: %v\n", err)
				return errutil.ErrFailedPrecondition
			}
		}
	}

	return nil
}

func (g *gitGateway) Head() (string, error) {
	ref, err := g.repo.Head()
	if err != nil {
		return "", errutil.ErrFailedPrecondition
	}
	return ref.Hash().String(), nil
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

func (g *gitGateway) Diff(fromHash string, toHash string) ([]string, error) {
	from := plumbing.NewHash(fromHash)
	to := plumbing.NewHash(toHash)

	fromCommit, err := g.repo.CommitObject(from)
	if err != nil {
		return []string{}, errutil.ErrFailedPrecondition
	}

	toCommit, err := g.repo.CommitObject(to)
	if err != nil {
		return []string{}, errutil.ErrFailedPrecondition
	}

	fromTree, err := fromCommit.Tree()
	if err != nil {
		return []string{}, errutil.ErrFailedPrecondition
	}

	toTree, err := toCommit.Tree()
	if err != nil {
		return []string{}, errutil.ErrFailedPrecondition
	}

	diff, err := object.DiffTree(fromTree, toTree)
	if err != nil {
		return []string{}, errutil.ErrFailedPrecondition
	}

	paths := []string{}
	for _, d := range diff {
		action, actionErr := d.Action()
		if actionErr != nil {
			return []string{}, errutil.ErrFailedPrecondition
		}

		if action == merkletrie.Delete || action == merkletrie.Modify {
			paths = append(paths, d.From.Name)
		}

		if action == merkletrie.Insert {
			paths = append(paths, d.To.Name)
		}
	}

	return paths, nil
}

func (g *gitGateway) Log() error {
	ref, refErr := g.repo.Head()
	if refErr != nil {
		return errutil.ErrFailedPrecondition
	}

	cIter, logErr := g.repo.Log(&git.LogOptions{From: ref.Hash()})
	if logErr != nil {
		return errutil.ErrFailedPrecondition
	}

	err := cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})

	if err != nil {
		return errutil.ErrUnknown
	}

	return nil

}
