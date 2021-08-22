package gitopscontroller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/gitopscontroller"
)

func TestGitOpsController_SetupUserspace(t *testing.T) {
	err := gitopscontroller.SetupUserspace()
	assert.Nil(t, err)
}
