package repositoryobserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	testSuites := []*struct {
		repoURL       string
		currentCommit string
		expect        string
		expectNil     bool
		branch        string
	}{
		{
			repoURL:       "https://github.com/thecodeisalreadydeployed/fixture-monorepo",
			currentCommit: "37e8e4d20d889924780f2373453a246591b6b11a", // feat(nx): init Nx workspace
			expect:        "5da29979c5ef986dc8ec6aa603e0862310abc96e", // build(dev-deps): @nrwl/next
			branch:        "main",
		},
		{
			repoURL:       "https://github.com/thecodeisalreadydeployed/fixture-monorepo",
			currentCommit: "5da29979c5ef986dc8ec6aa603e0862310abc96e",
			expect:        "",
			branch:        "main",
		},
		{
			repoURL:       "https://github.com/thecodeisalreadydeployed/fixture-nest",
			currentCommit: "62139be31792ab4a43c00eadcc8af6cadd90ee66", // feat: init NestJS project
			expect:        "14bc77fc515e6d66b8d9c15126ee49ca55faf879", // chore(app): Hello World -> fixture-nest
			branch:        "main",
		},
	}

	for _, test := range testSuites {
		commit := check(test.repoURL, test.branch, test.currentCommit)
		assert.Equal(t, test.expect, commit)
	}
}
