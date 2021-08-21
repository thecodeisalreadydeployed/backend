package gitinteractor

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	it "github.com/thecodeisalreadydeployed/gitinteractor"
)

var privateKey = flag.String("private-key", "", "")

func TestGitInteractor_Clone(t *testing.T) {
	git := it.NewGitInteractorSSH("ssh://codedeploy@localhost:2222/__w/repos/userspace", *privateKey)
	log := git.Log()
	assert.Equal(t, 1, len(log))
}
