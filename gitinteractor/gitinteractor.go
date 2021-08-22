package gitinteractor

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/logger"
	"go.uber.org/zap"
	gossh "golang.org/x/crypto/ssh"
)

type GitInteractor struct {
	repository *git.Repository
	storage    storage.Storer
	fs         billy.Filesystem
}

func NewGitInteractor() GitInteractor {
	it := GitInteractor{}
	it.storage = memory.NewStorage()
	it.fs = memfs.New()
	it.Init()
	return it
}

func NewGitInteractorSSH(url string, privateKey string) GitInteractor {
	signer, parsePrivateKeyErr := gossh.ParsePrivateKey([]byte(privateKey))
	if parsePrivateKeyErr != nil {
		panic(parsePrivateKeyErr)
	}

	auth := &ssh.PublicKeys{
		User:   "codedeploy",
		Signer: signer,
	}

	auth.HostKeyCallback = gossh.InsecureIgnoreHostKey()

	it := GitInteractor{}

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  url,
		Auth: auth,
	})
	if err != nil {
		spew.Dump(err)
		panic(err)
	}

	it.repository = r
	return it
}

func InitRepository(path string) error {
	_, err := git.PlainInit(path, false)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to create Git repository at path: %s", path), zap.String("package", "gitinteractor"))
		return err
	}
	return nil
}

func (it *GitInteractor) Init() {
	repo, err := git.Init(it.storage, it.fs)
	if err != nil {
		panic(err)
	}
	it.repository = repo
}

func (it *GitInteractor) Add(filePath string) {
	w, err := it.repository.Worktree()
	if err != nil {
		panic(err)
	}
	w.Add(filePath)
}

func (it *GitInteractor) Commit(message string) {
	w, err := it.repository.Worktree()
	if err != nil {
		panic(err)
	}
	w.Commit(message, &git.CommitOptions{Author: config.DefaultGitSignature()})
}

func (it *GitInteractor) WriteFile(path string, name string, data []byte) {
	w, err := it.repository.Worktree()
	if err != nil {
		panic(err)
	}

	fs := w.Filesystem

	err = fs.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}

	util.WriteFile(fs, fs.Join(path, name), data, 0644)
}

func (it *GitInteractor) Log() []string {
	messages := []string{}
	r := it.repository
	ref, err := r.Head()
	if err != nil {
		panic(err)
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		messages = append(messages, c.Message)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return messages
}
