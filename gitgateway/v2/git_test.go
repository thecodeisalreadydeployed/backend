package gitgateway

import (
	"os"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func temporalDir() (path string, clean func()) {
	fs := osfs.New(os.TempDir())
	path, err := util.TempDir(fs, "", "")
	if err != nil {
		panic(err)
	}
	return fs.Join(fs.Root(), path), func() {
		err := util.RemoveAll(fs, path)
		if err != nil {
			panic(err)
		}
	}
}

func initRepository() (path string, clean func()) {
	dir, clean := temporalDir()
	_, initErr := git.PlainInit(dir, false)
	if initErr != nil {
		panic(initErr)
	}
	return dir, clean
}

func TestGitGateway(t *testing.T) {
	path, clean := initRepository()
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
}
