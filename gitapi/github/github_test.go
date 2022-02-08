package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestListBranches(t *testing.T) {
	logger := zaptest.NewLogger(t)
	gitHubAPI := NewGitHubAPI(logger, "octocat", "Hello-World")
	branches, err := gitHubAPI.GetBranches()
	assert.Nil(t, err)
	assert.ElementsMatch(t, [3]string{"master", "test", "octocat-patch-1"}, branches)
}

func TestListFiles(t *testing.T) {
	logger := zaptest.NewLogger(t)
	gitHubAPI := NewGitHubAPI(logger, "octocat", "Hello-World")
	files, err := gitHubAPI.GetFiles("master")
	assert.Nil(t, err)
	assert.ElementsMatch(t, [1]string{"README"}, files)
}

func TestGetRaw(t *testing.T) {
	logger := zaptest.NewLogger(t)
	gitHubAPI := NewGitHubAPI(logger, "octocat", "Hello-World")
	raw, err := gitHubAPI.GetRaw("master", "README")
	assert.Nil(t, err)
	assert.Equal(t, "Hello World!\n", raw)
}
