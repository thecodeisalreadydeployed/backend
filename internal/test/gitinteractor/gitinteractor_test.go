package gitinteractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	it "github.com/thecodeisalreadydeployed/gitinteractor"
)

func TestGitInteractor_Clone(t *testing.T) {
	git := it.NewGitInteractor()
	_ = git
	assert.True(t, true)
}
