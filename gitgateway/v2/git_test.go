package gitgateway

import (
	"bou.ke/monkey"
	"github.com/thecodeisalreadydeployed/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGitGateway(t *testing.T) {
	path, clean := InitRepository()
	defer clean()

	git, err := NewGitGatewayLocal(path)
	assert.Nil(t, err)

	err = git.WriteFile(".thecodeisalreadydeployed", "")
	assert.Nil(t, err)

	_, err = git.Commit([]string{".thecodeisalreadydeployed"}, "Initial commit")
	assert.Nil(t, err)

	err = git.Checkout("deploy")
	assert.Nil(t, err)

	err = git.WriteFile(".thecodeisalreadydeployed", "A")
	assert.Nil(t, err)

	err = git.WriteFile(".thecodeisalreadydeployed", "B")
	assert.Nil(t, err)

	data, err := git.OpenFile(".thecodeisalreadydeployed")
	assert.Nil(t, err)
	assert.Equal(t, "B", data)

	b, err := git.Commit([]string{".thecodeisalreadydeployed"}, "B")
	assert.Nil(t, err)
	assert.NotEmpty(t, b)

	err = git.WriteFile(".thecodeisalreadydeployed", "C")
	assert.Nil(t, err)

	data, err = git.OpenFile(".thecodeisalreadydeployed")
	assert.Nil(t, err)
	assert.Equal(t, "C", data)

	c, err := git.Commit([]string{".thecodeisalreadydeployed"}, "C")
	assert.Nil(t, err)
	assert.NotEmpty(t, c)
	assert.NotEqual(t, b, c)

	diff, err := git.Diff(b, c)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(diff))
	assert.Equal(t, ".thecodeisalreadydeployed", diff[0])

	err = git.Log()
	assert.Nil(t, err)
}

func TestCommitInterval(t *testing.T) {
	defer monkey.UnpatchAll()

	path, clean := InitRepository()
	defer clean()

	git, err := NewGitGatewayLocal(path)
	assert.Nil(t, err)

	commit(time.Unix(0, 0), git, t)
	commit(time.Unix(50, 0), git, t)
	commit(time.Unix(150, 0), git, t)

	duration, err := git.CommitInterval()
	assert.Nil(t, err)
	assert.Equal(t, 75*time.Second, duration)
}

func commit(commitTime time.Time, git GitGateway, t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return commitTime
	})

	err := git.WriteFile(".thecodeisalreadydeployed", "data")
	assert.Nil(t, err)

	_, err = git.Commit([]string{".thecodeisalreadydeployed"}, "This is a commit.")
	assert.Nil(t, err)
}

func TestGetBranches(t *testing.T) {
	git, err := NewGitGatewayRemote("https://github.com/octocat/Hello-World")
	assert.Nil(t, err)

	branches, err := git.GetBranches()
	assert.Nil(t, err)
	assert.ElementsMatch(t, branches, [3]string{"master", "test", "octocat-patch-1"})
}

func TestGetFiles(t *testing.T) {
	git, err := NewGitGatewayRemote("https://github.com/octocat/Hello-World")
	assert.Nil(t, err)

	files, err := git.GetFiles("test")
	assert.Nil(t, err)
	assert.ElementsMatch(t, files, [2]string{"CONTRIBUTING.md", "README"})
}

func TestGetRaw(t *testing.T) {
	git, err := NewGitGatewayRemote("https://github.com/octocat/Hello-World")
	assert.Nil(t, err)

	raw, err := git.GetRaw("master", "README")
	assert.Nil(t, err)
	assert.Equal(t, raw, "Hello World!\n")
}

func TestInfo(t *testing.T) {
	gs, err := Info("https://github.com/octocat/Hello-World", "master")
	assert.Nil(t, err)
	assert.Equal(t, model.GitSource{
		CommitSHA:        "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		CommitMessage:    "Merge pull request #6 from Spaceghost/patch-1\n\nNew line at end of file.",
		CommitAuthorName: "The Octocat",
		RepositoryURL:    "https://github.com/octocat/Hello-World",
		Branch:           "master",
	}, gs)
}
