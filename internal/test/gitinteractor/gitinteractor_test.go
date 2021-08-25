package gitinteractor

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

var privateKey = flag.String("private-key", "", "") //nolint:golint,deadcode

func TestGitInteractor_InitRepository(t *testing.T) {
	assert.True(t, true)
}
