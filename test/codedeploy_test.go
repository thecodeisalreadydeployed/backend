package test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/model"
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
	if os.Getenv("GITHUB_WORKFLOW") == "test: kind" {
		time.Sleep(30 * time.Second)
	}

	expect.POST("/apps").
		WithForm(dto.CreateAppRequest{
			ProjectID:     projectID,
			Name:          appName,
			RepositoryURL: "https://github.com/thecodeisalreadydeployed/fixture-nest",
			BuildScript:   fixtureNest,
			Branch:        "main",
		}).Expect().Status(http.StatusOK)

	apps := expect.GET("/projects/" + projectID + "/apps").
		Expect().
		Status(http.StatusOK).
		JSON()

	apps.Array().Length().Equal(1)

	appID := apps.Array().Element(0).Object().Value("id").String().Raw()
	assert.NotEmpty(t, appID)

	expect.POST(fmt.Sprintf("/apps/%s/observable/disable", appID)).Expect().Status(http.StatusOK)

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
		Expect().Status(http.StatusOK).JSON().Null()

	actualProject := expect.GET("/projects/search").WithQuery("name", projectName).
		Expect().Status(http.StatusOK).JSON().
		Array().Element(0).Object()

	assert.Equal(t, projectID, actualProject.Value("id").String().Raw())
	assert.Equal(t, projectName, actualProject.Value("name").String().Raw())

	actualApp := expect.GET("/apps/search").WithQuery("name", appName).
		Expect().Status(http.StatusOK).JSON().
		Array().Element(0).Object()

	assert.Equal(t, appID, actualApp.Value("id").String().Raw())
	assert.Equal(t, appName, actualApp.Value("name").String().Raw())
	deployment := expect.POST(fmt.Sprintf("/apps/%s/deployments", appID)).
		Expect().Status(http.StatusOK).
		JSON()

	deployment.Object().
		ContainsKey("state").
		ValueEqual("state", model.DeploymentStateQueueing)

	if os.Getenv("GITHUB_WORKFLOW") == "test: kind" {
		deploymentID := deployment.Object().Value("id").String().Raw()

		timeLimit := time.Now().Add(1 * time.Minute)
		for {
			if time.Now().After(timeLimit) {
				t.Fatal("didn't see result in time")
			}

			deployment := expect.GET(fmt.Sprintf("/deployments/%s", deploymentID)).Expect().Status(http.StatusOK).JSON()
			deploymentState := deployment.Object().Value("state").String().Raw()
			if deploymentState != string(model.DeploymentStateBuilding) {
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}

		timeLimit = time.Now().Add(30 * time.Second)
		for {
			if time.Now().After(timeLimit) {
				t.Fatal("didn't see result in time")
			}

			events := expect.GET("/deployments/" + deploymentID + "/events").
				Expect().
				Status(http.StatusOK).
				JSON().Array().Raw()

			if len(events) == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}

		timeLimit = time.Now().Add(5 * time.Minute)
		for {
			if time.Now().After(timeLimit) {
				t.Fatal("didn't see result in time")
			}

			deployment := expect.GET(fmt.Sprintf("/deployments/%s", deploymentID)).Expect().Status(http.StatusOK).JSON()
			deploymentState := deployment.Object().Value("state").String().Raw()
			if deploymentState != string(model.DeploymentStateBuildSucceeded) {
				if deploymentState == string(model.DeploymentStateCommitted) {
					break
				}
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}

		timeLimit = time.Now().Add(1 * time.Minute)
		for {
			if time.Now().After(timeLimit) {
				t.Fatal("didn't see result in time")
			}

			deployment := expect.GET(fmt.Sprintf("/deployments/%s", deploymentID)).Expect().Status(http.StatusOK).JSON()
			deploymentState := deployment.Object().Value("state").String().Raw()
			if deploymentState != string(model.DeploymentStateCommitted) {
				if deploymentState == string(model.DeploymentStateReady) {
					break
				}
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}

		timeLimit = time.Now().Add(15 * time.Minute)
		for {
			if time.Now().After(timeLimit) {
				t.Fatal("didn't see result in time")
			}

			deployment := expect.GET(fmt.Sprintf("/deployments/%s", deploymentID)).Expect().Status(http.StatusOK).JSON()
			deploymentState := deployment.Object().Value("state").String().Raw()
			if deploymentState != string(model.DeploymentStateReady) {
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}

		timeLimit = time.Now().Add(15 * time.Minute)
		for {
			if time.Now().After(timeLimit) {
				t.Fatal("didn't see result in time")
			}

			appStatus := expect.GET(fmt.Sprintf("/apps/%s/status", appID)).Expect().Status(http.StatusOK).JSON()
			appStatusDeploymentID := appStatus.Object().Value("deploymentID").String().Raw()
			if appStatusDeploymentID != deploymentID {
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}
	}

	// expect.DELETE(fmt.Sprintf("/apps/%s", appID)).
	// 	Expect().Status(http.StatusOK)

	// expect.DELETE(fmt.Sprintf("/projects/%s", projectID)).
	// 	Expect().Status(http.StatusOK)

	// expect.GET(fmt.Sprintf("/projects/%s", projectID)).
	// 	Expect().Status(http.StatusNotFound)

	// expect.GET(fmt.Sprintf("/apps/%s", appID)).
	// 	Expect().Status(http.StatusNotFound)
}

func TestPresetIntegration(t *testing.T) {
	presetName := "best-preset"

	expect := Setup(t)

	expect.POST("/presets").
		WithForm(dto.CreatePresetRequest{Name: presetName, Template: "RUN echo hello"}).
		Expect().
		Status(http.StatusOK)

	presets := expect.GET("/presets/list").
		Expect().
		Status(http.StatusOK).
		JSON()

	presets.Array().Length().Equal(5)

	presetID := presets.Array().Element(4).Object().Value("id").String().Raw()
	assert.NotEmpty(t, presetID)

	expect.GET(fmt.Sprintf("/presets/%s", presetID)).
		Expect().Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("name").
		ValueEqual("name", presetName)

	actualPreset := expect.GET("/presets/search").WithQuery("name", presetName).
		Expect().
		Status(http.StatusOK).
		JSON().
		Array().
		Element(0)

	assert.Equal(t, presetID, actualPreset.Object().Value("id").String().Raw())
	assert.Equal(t, presetName, actualPreset.Object().Value("name").String().Raw())

	expect.DELETE(fmt.Sprintf("/presets/%s", presetID)).
		Expect().Status(http.StatusOK)

	expect.GET(fmt.Sprintf("/presets/%s", presetID)).
		Expect().Status(http.StatusNotFound)
}

func TestGitHubApiIntegration(t *testing.T) {
	expect := Setup(t)

	branches := expect.POST("/gitapi/branches").
		WithForm(dto.GetBranchesRequest{URL: "https://github.com/octocat/Hello-World"}).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	assert.ElementsMatch(t, (*branches).Raw(), [3]string{"master", "test", "octocat-patch-1"})

	files := expect.POST("/gitapi/files").
		WithForm(dto.GetFilesRequest{
			URL:    "https://github.com/octocat/Hello-World",
			Branch: "test",
		}).Expect().Status(http.StatusOK).JSON().Array()

	assert.ElementsMatch(t, (*files).Raw(), [2]string{"CONTRIBUTING.md", "README"})

	raw := expect.POST("/gitapi/raw").
		WithForm(dto.GetRawRequest{
			URL:    "https://github.com/octocat/Hello-World",
			Branch: "octocat-patch-1",
			Path:   "README",
		}).Expect().Status(http.StatusOK).Body().Raw()

	assert.Equal(t, "Hello world!\n", raw)
}

func TestGitLabApiIntegration(t *testing.T) {
	expect := Setup(t)

	branches := expect.POST("/gitapi/branches").
		WithForm(dto.GetBranchesRequest{URL: "https://gitlab.com/gitlab-examples/docker"}).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	assert.ElementsMatch(t, (*branches).Raw(), [1]string{"master"})

	files := expect.POST("/gitapi/files").
		WithForm(dto.GetFilesRequest{
			URL:    "https://gitlab.com/gitlab-examples/docker",
			Branch: "master",
		}).Expect().Status(http.StatusOK).JSON().Array()

	assert.ElementsMatch(t, (*files).Raw(), [3]string{".gitlab-ci.yml", "Dockerfile", "README.md"})

	raw := expect.POST("/gitapi/raw").
		WithForm(dto.GetRawRequest{
			URL:    "https://gitlab.com/gitlab-examples/docker",
			Branch: "master",
			Path:   "Dockerfile",
		}).Expect().Status(http.StatusOK).Body().Raw()

	expected := "FROM alpine:latest\nRUN apk add -U git\n"

	assert.Equal(t, expected, raw)
}
