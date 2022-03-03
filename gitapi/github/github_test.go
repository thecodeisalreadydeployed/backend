package github

import (
	"github.com/thecodeisalreadydeployed/model"
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

func TestFillGitSource(t *testing.T) {
	logger := zaptest.NewLogger(t)
	gitHubAPI := NewGitHubAPI(logger, "octocat", "Hello-World")
	gs, err := gitHubAPI.FillGitSource(&model.GitSource{
		RepositoryURL: "https://github.com/octocat/Hello-World.git",
		Branch:        "master",
	})
	assert.Nil(t, err)
	assert.Equal(t, &model.GitSource{
		CommitSHA:        "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		CommitMessage:    "Merge pull request #6 from Spaceghost/patch-1\n\nNew line at end of file.",
		CommitAuthorName: "The Octocat",
		RepositoryURL:    "https://github.com/octocat/Hello-World.git",
		Branch:           "master",
	}, gs)
}
