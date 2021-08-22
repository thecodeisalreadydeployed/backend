package gitopscontroller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/gitopscontroller"
)

func TestGitOpsController_Log(t *testing.T) {
	err := gitopscontroller.SetupUserspace()
	assert.Nil(t, err)

	controller := gitopscontroller.GitOpsController{}
	userspaceRepository := controller.Userspace
	assert.Equal(t, 0, len(userspaceRepository.Log()))

}
