package gitopscontroller_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/constant"
	"github.com/thecodeisalreadydeployed/gitopscontroller"
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

func TestGitOpsController(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("GITHUB_WORKFLOW") == "test: unit" {
		dir, clean := temporalDir()
		viper.Set(constant.UserspaceRepository, dir)
		gitopscontroller.SetupUserspace()
		controller := gitopscontroller.NewGitOpsController()
		err := controller.SetupApp("prj-test", "app-test")
		assert.NoError(t, err)

		_, err = os.Stat(filepath.Join(dir, "prj-test", "app-test", "deployment.yml"))
		assert.NoError(t, err)

		_, err = os.Stat(filepath.Join(dir, "prj-test", "app-test", "service.yml"))
		assert.NoError(t, err)

		_, err = os.Stat(filepath.Join(dir, "prj-test", "app-test", "kustomization.yml"))
		assert.NoError(t, err)

		err = controller.SetContainerImage("prj-test", "app-test", "ghcr.io/thecodeisalreadydeployed/imagebuilder:latest")
		assert.NoError(t, err)

		clean()
	}
}
