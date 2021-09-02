package gitopscontroller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/gitopscontroller"
)

func TestGitOpsController_Log(t *testing.T) {

	gitopscontroller.Init()
	controller := gitopscontroller.GetController()

	err := controller.SetupUserspace()
	assert.Nil(t, err)

	controller.Write("/.thecodeisalreadydeployed", "")
}
