package gitapi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListBranches(t *testing.T) {
	branches, err := ListBranches("https://github.com/octocat/Hello-World")
	assert.Nil(t, err)
	assert.ElementsMatch(t, branches, [3]string{"master", "test", "octocat-patch-1"})
}

func TestListFiles(t *testing.T) {
	files, err := ListFiles("https://github.com/octocat/Hello-World", "master")
	assert.Nil(t, err)
	assert.ElementsMatch(t, files, [1]string{"README"})
}

func TestGetRaw(t *testing.T) {
	raw, err := GetRaw("https://github.com/octocat/Hello-World", "master", "README")
	assert.Nil(t, err)
	assert.Equal(t, raw, "Hello World!\n")
}
