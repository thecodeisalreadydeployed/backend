package gitgateway

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
	"go.uber.org/zap"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/thecodeisalreadydeployed/config"
	gossh "golang.org/x/crypto/ssh"
)

type GitGateway struct {
	repository *git.Repository
}

func NewGitGateway(path string) GitGateway {
	it := GitGateway{}
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}
	it.repository = repo
	return it
}

func NewGitGatewaySSH(url string, privateKey string) GitGateway {
	signer, parsePrivateKeyErr := gossh.ParsePrivateKey([]byte(privateKey))
	if parsePrivateKeyErr != nil {
		panic(parsePrivateKeyErr)
	}

	auth := &ssh.PublicKeys{
		User:   "codedeploy",
		Signer: signer,
	}

	auth.HostKeyCallback = gossh.InsecureIgnoreHostKey()

	it := GitGateway{}

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
		zap.L().Error(fmt.Sprintf("Failed to create Git repository at path: %s", path))
		return err
	}
	return nil
}

func (it *GitGateway) Add(filePath string) {
	w, err := it.repository.Worktree()
	if err != nil {
		panic(err)
	}
	_, err = w.Add(filePath)
	if err != nil {
		panic(err)
	}
}

func (it *GitGateway) Commit(message string) {
	w, err := it.repository.Worktree()
	if err != nil {
		panic(err)
	}
	_, err = w.Commit(message, &git.CommitOptions{Author: config.DefaultGitSignature()})
	if err != nil {
		panic(err)
	}
}

func (it *GitGateway) WriteFile(path string, name string, data []byte) {
	w, err := it.repository.Worktree()
	if err != nil {
		panic(err)
	}

	fs := w.Filesystem

	err = fs.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}

	err = util.WriteFile(fs, fs.Join(path, name), data, 0644)
	if err != nil {
		panic(err)
	}
}

func (it *GitGateway) Log() []string {
	messages := []string{}
	r := it.repository
	ref, err := r.Head()
	if err != nil {
		panic(err)
	}
	cIter, logErr := r.Log(&git.LogOptions{From: ref.Hash()})
	if logErr != nil {
		panic(err)
	}
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

func (it *GitGateway) CreateBranch(name string) {
	err := it.repository.CreateBranch(&gitconfig.Branch{Name: name})
	if err != nil {
		panic(err)
	}
}

func (it *GitGateway) Checkout(name string) {
	w, err := it.repository.Worktree()
	if err != nil {
		panic(err)
	}

	branch := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", name))
	err = w.Checkout(&git.CheckoutOptions{
		Hash:   plumbing.Hash{},
		Branch: branch,
		Create: false,
		Force:  false,
		Keep:   false,
	})
	if err != nil {
		panic(err)
	}
}

func (it *GitGateway) GetCommitSHA() string {
	ref, err := it.repository.Head()
	if err != nil {
		panic(err)
	}

	return ref.Hash().String()
}

func (it *GitGateway) GetCommit(hash string) *object.Commit {
	commit, err := it.repository.CommitObject(plumbing.NewHash(hash))
	if err != nil {
		panic(err)
	}

	return commit
}

func (it *GitGateway) GetCurrentCommit() *object.Commit {
	ref, err := it.repository.Head()
	if err != nil {
		panic(err)
	}

	commit, err := it.repository.CommitObject(ref.Hash())
	if err != nil {
		panic(err)
	}

	return commit
}

func diff(old *object.Commit, current *object.Commit) []string {
	patch, err := old.Patch(current)
	if err != nil {
		panic(err)
	}

	stats := patch.Stats()
	var files []string

	for _, stat := range stats {
		files = append(files, stat.Name)
	}
	return files
}

func HasProperDiff(old *object.Commit, current *object.Commit) bool {
	files := diff(old, current)
	for _, file := range files {
		for _, ignore := range gitignore {
			if file != ignore {
				return true
			}
		}
	}
	return false
}
