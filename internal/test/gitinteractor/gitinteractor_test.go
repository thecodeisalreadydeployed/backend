package gitinteractor

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	it "github.com/thecodeisalreadydeployed/gitinteractor"
)

var privateKey = pflag.String("private-key", "", "")

func TestGitInteractor_Clone(t *testing.T) {
	pflag.Parse()
	git := it.NewGitInteractorSSH("git@localhost:2222/__w/repos/userspace", *privateKey)
	log := git.Log()
	assert.Equal(t, 1, len(log))
}
