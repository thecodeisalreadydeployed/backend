package gitgateway

import (
	"bou.ke/monkey"
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

func TestCommitDuration(t *testing.T) {
	defer monkey.UnpatchAll()

	path, clean := InitRepository()
	defer clean()

	git, err := NewGitGatewayLocal(path)
	assert.Nil(t, err)

	commit(time.Unix(0, 0), git, t)
	commit(time.Unix(50, 0), git, t)
	commit(time.Unix(150, 0), git, t)

	duration, err := git.CommitDuration()
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
