package gitopscontroller_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/gitopscontroller"
)

func TestGitOpsController(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("GITHUB_WORKFLOW") == "test: unit" {
		controller := gitopscontroller.NewGitOpsController()
		err := controller.SetupApp("prj-test", "app-test")
		assert.NoError(t, err)
	}
}
