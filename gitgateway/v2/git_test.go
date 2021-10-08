package gitgateway

import (
	"os"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
)

func temporalDir() (path string, clean func()) {
	fs := osfs.New(os.TempDir())
	path, err := util.TempDir(fs, "", "")
	if err != nil {
		panic(err)
	}
	return fs.Join(fs.Root(), path), func() {
		util.RemoveAll(fs, path)
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
