package gitapi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListBranches(t *testing.T) {
	branches := ListBranches("https://github.com/octocat/Hello-World")
	assert.ElementsMatch(t, branches, [3]string{"master", "test", "octocat-patch-1"})
}

func TestListFiles(t *testing.T) {
	files := ListFiles("https://github.com/octocat/Hello-World", "master")
	assert.ElementsMatch(t, files, [1]string{"README.md"})
}

func TestGetRaw(t *testing.T) {
	raw := GetRaw("https://github.com/octocat/Hello-World", "master", "README.md")
	assert.Equal(t, raw, "Hello World!")
}
