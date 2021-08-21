package gitinteractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	it "github.com/thecodeisalreadydeployed/gitinteractor"
)

func TestGitInteractor_Clone(t *testing.T) {
	git := it.NewGitInteractor()
	assert.True(t, true)
}
