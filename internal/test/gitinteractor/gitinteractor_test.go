package gitinteractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	it "github.com/thecodeisalreadydeployed/gitinteractor"
)

func TestGitInteractor_Clone(t *testing.T) {
	git := it.NewGitInteractorSSH("git@localhost:2222/__w/repos/userspace")
	log := git.Log()
	assert.Equal(t, 1, len(log))
}
