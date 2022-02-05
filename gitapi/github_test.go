package gitapi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListBranches(t *testing.T) {
	branches, err := GetBranches("https://github.com/octocat/Hello-World")
	assert.Nil(t, err)
	assert.ElementsMatch(t, [3]string{"master", "test", "octocat-patch-1"}, branches)
}

func TestListFiles(t *testing.T) {
	files, err := GetFiles("https://github.com/octocat/Hello-World", "master")
	assert.Nil(t, err)
	assert.ElementsMatch(t, [1]string{"README"}, files)
}

func TestGetRaw(t *testing.T) {
	raw, err := GetRaw("https://github.com/octocat/Hello-World", "master", "README")
	assert.Nil(t, err)
	assert.Equal(t, "Hello World!\n", raw)
}
