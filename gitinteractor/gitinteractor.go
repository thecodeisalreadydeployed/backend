package gitinteractor

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/thecodeisalreadydeployed/config"
)

type GitInteractor struct {
	Repository *git.Repository
}

func NewGitInteractor() GitInteractor {
	it := GitInteractor{}
	it.Init()
	return it
}

func (it *GitInteractor) Init() {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	if err != nil {
		panic(err)
	}
	it.Repository = repo
}

func (it *GitInteractor) Commit(message string) {
	w, err := it.Repository.Worktree()
	if err != nil {
		panic(err)
	}
	w.Commit(message, &git.CommitOptions{Author: config.DefaultGitSignature()})
}
