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
	controller.Init()
	userspaceRepository := controller.Userspace
	userspaceRepository.WriteFile("/", ".thecodeisalreadydeployed", []byte(""))
	userspaceRepository.Add("/.thecodeisalreadydeployed")
	userspaceRepository.Commit(".thecodeisalreadydeployed: init")
	assert.Equal(t, 1, len(userspaceRepository.Log()))
}
