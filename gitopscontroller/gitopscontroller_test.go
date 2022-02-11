package gitopscontroller_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/constant"
	"github.com/thecodeisalreadydeployed/gitopscontroller"
	"github.com/thecodeisalreadydeployed/gitopscontroller/argocd"
	"go.uber.org/zap/zaptest"
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
		viper.Set(constant.ARGOCD_SERVER_HOST, "http://localhost")
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		defaultHTTPTransport := argocd.HTTPTransport
		defer func() { argocd.HTTPTransport = defaultHTTPTransport }()
		argocd.HTTPTransport = httpmock.DefaultTransport

		httpmock.RegisterResponder(
			"GET",
			"http://localhost/api/v1/applications?name=userspace&refresh=true",
			httpmock.NewStringResponder(200, ""),
		)

		httpmock.RegisterResponder(
			"POST",
			"http://localhost/api/v1/applications",
			httpmock.NewStringResponder(200, ""),
		)

		httpmock.RegisterResponder(
			"POST",
			"http://localhost/api/v1/applications/userspace/sync",
			httpmock.NewStringResponder(200, ""),
		)

		dir, clean := temporalDir()
		viper.Set(constant.USERSPACE_REPOSITORY, dir)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := zaptest.NewLogger(t)
		controller := gitopscontroller.NewGitOpsController(logger)

		err := controller.SetupProject("prj-test")
		assert.NoError(t, err)

		contents, err := os.ReadFile(filepath.Join(dir, "kustomization.yml"))
		assert.NoError(t, err)
		fmt.Printf("contents (/kustomization.yml): %v\n", contents)

		err = controller.SetupApp("prj-test", "app-test")
		assert.NoError(t, err)

		contents, err = os.ReadFile(filepath.Join(dir, "prj-test", "kustomization.yml"))
		assert.NoError(t, err)
		fmt.Printf("contents (/prj-test/kustomization.yml): %v\n", contents)

		_, err = os.Stat(filepath.Join(dir, "prj-test", "app-test", "deployment.yml"))
		assert.NoError(t, err)

		contents, err = os.ReadFile(filepath.Join(dir, "prj-test", "app-test", "deployment.yml"))
		assert.NoError(t, err)
		fmt.Printf("contents (/prj-test/app-test/deployment.yml): %v\n", contents)

		_, err = os.Stat(filepath.Join(dir, "prj-test", "app-test", "service.yml"))
		assert.NoError(t, err)

		contents, err = os.ReadFile(filepath.Join(dir, "prj-test", "app-test", "service.yml"))
		assert.NoError(t, err)
		fmt.Printf("contents (/prj-test/app-test/service.yml): %v\n", contents)

		_, err = os.Stat(filepath.Join(dir, "prj-test", "app-test", "kustomization.yml"))
		assert.NoError(t, err)

		contents, err = os.ReadFile(filepath.Join(dir, "prj-test", "app-test", "kustomization.yml"))
		assert.NoError(t, err)
		fmt.Printf("contents (/prj-test/app-test/kustomization.yml): %v\n", contents)

		err = controller.SetContainerImage("prj-test", "app-test", "ghcr.io/thecodeisalreadydeployed/imagebuilder:latest")
		assert.NoError(t, err)

		contents, err = os.ReadFile(filepath.Join(dir, "prj-test", "app-test", "kustomization.yml"))
		assert.NoError(t, err)
		fmt.Printf("contents (/prj-test/app-test/kustomization.yml): %v\n", contents)

		clean()
	}
}
