package test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
)

func testFixture(t *testing.T, fixtureRepo string) {
	repoURL, urlParseErr := url.Parse(fixtureRepo)
	assert.Nil(t, urlParseErr)

	repo := strings.TrimLeft(repoURL.Path, "/")

	projectName := "project/" + repo
	appName := repo

	expect := Setup(t)

	expect.POST("/project").
		WithForm(dto.CreateProjectRequest{Name: projectName}).
		Expect().
		Status(http.StatusOK)
}

func TestFixture(t *testing.T) {
	testSuites := []*struct {
		repoURL          string
		buildScript      string
		installCommand   string
		buildCommand     string
		outputDirectory  string
		startCommand     string
		currentCommitSHA string
	}{
		{
			repoURL:          "https://github.com/thecodeisalreadydeployed/fixture-nest",
			currentCommitSHA: "62139be31792ab4a43c00eadcc8af6cadd90ee66", // v1
		},
		{
			repoURL: "https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		},
	}

	for _, tt := range testSuites {
		t.Run(tt.repoURL, func(t *testing.T) {
			testFixture(t, tt.repoURL)
		})
	}
}
