package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
)

func Setup(t *testing.T) *httpexpect.Expect {

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:3000",
		Reporter: httpexpect.NewAssertReporter(t),
	})

	return e
}

func TestHealth(t *testing.T) {
	expect := Setup(t)

	expect.GET("/health").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("ok").
		ValueEqual("ok", "true")
}

func TestIntegration(t *testing.T) {
	projectName := "deploys-dev"
	appName := "fixture-nest"
	fixtureNest := `
FROM node:14-alpine as build-env
ADD . /app
WORKDIR /app
RUN yarn install --frozen-lockfile
RUN yarn build

FROM node:14-alpine
WORKDIR /app
COPY --from=build-env /app/package.json /app/yarn.lock ./
COPY --from=build-env /app/node_modules ./node_modules
COPY --from=build-env /app/dist ./
CMD node main
`

	expect := Setup(t)

	expect.POST("/projects").
		WithForm(dto.CreateProjectRequest{Name: projectName}).
		Expect().
		Status(http.StatusOK)

	projects := expect.GET("/projects/list").
		Expect().
		Status(http.StatusOK).
		JSON()

	projects.Array().Length().Equal(1)

	projectID := projects.Array().Element(0).Object().Value("id").String().Raw()

	assert.NotEmpty(t, projectID)

	expect.POST("/apps").
		WithForm(dto.CreateAppRequest{
			ProjectID:       projectID,
			Name:            appName,
			RepositoryURL:   "https://github.com/thecodeisalreadydeployed/fixture-nest.git",
			BuildScript:     fixtureNest,
			InstallCommand:  "yarn install --frozen-lockfile",
			BuildCommand:    "yarn build",
			OutputDirectory: "dist",
			StartCommand:    "node main",
			Branch:          "main",
		}).Expect().Status(http.StatusOK)

	apps := expect.GET("/projects/" + projectID + "/apps").
		Expect().
		Status(http.StatusOK).
		JSON()

	apps.Array().Length().Equal(1)

	appID := apps.Array().Element(0).Object().Value("id").String().Raw()

	assert.NotEmpty(t, appID)

	expect.GET(fmt.Sprintf("/projects/%s", projectID)).
		Expect().Status(http.StatusOK).JSON().
		Object().
		ContainsKey("name").ValueEqual("name", projectName)

	expect.GET(fmt.Sprintf("/apps/%s", appID)).
		Expect().Status(http.StatusOK).JSON().
		Object().
		ContainsKey("projectID").ValueEqual("projectID", projectID).
		ContainsKey("name").ValueEqual("name", appName)

	expect.GET(fmt.Sprintf("/apps/%s/deployments", appID)).
		Expect().Status(http.StatusOK).JSON()

	expect.POST(fmt.Sprintf("/apps/%s/deployments", appID)).
		Expect().Status(http.StatusOK).
		JSON().Object().
		ContainsKey("ok").
		ValueEqual("ok", "true")

	expect.DELETE(fmt.Sprintf("/apps/%s", appID)).
		Expect().Status(http.StatusOK)

	expect.DELETE(fmt.Sprintf("/projects/%s", projectID)).
		Expect().Status(http.StatusOK)

	expect.GET(fmt.Sprintf("/projects/%s", projectID)).
		Expect().Status(http.StatusNotFound)

	expect.GET(fmt.Sprintf("/apps/%s", appID)).
		Expect().Status(http.StatusNotFound)
}
