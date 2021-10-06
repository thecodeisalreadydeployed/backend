package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
)

func setup(t *testing.T) *httpexpect.Expect {

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:3000",
		Reporter: httpexpect.NewAssertReporter(t),
	})

	return e
}

func TestHealth(t *testing.T) {
	expect := setup(t)

	expect.GET("/health").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("ok").
		ValueEqual("ok", "true")
}

func TestIntegration(t *testing.T) {
	projectName := "Test Project"
	appName := "Test App"
	fake := "Fake Data"

	expect := setup(t)

	expect.POST("/project").
		WithForm(dto.CreateProjectRequest{Name: projectName}).
		Expect().
		Status(http.StatusOK)

	projects := expect.GET("/projects").
		Expect().
		Status(http.StatusOK).
		JSON()

	projects.Array().Length().Equal(1)

	projectID := projects.Array().Element(0).Object().Value("id").String().Raw()

	expect.POST("/app").
		WithForm(dto.CreateAppRequest{
			ProjectID:       projectID,
			Name:            appName,
			RepositoryURL:   fake,
			BuildScript:     fake,
			InstallCommand:  fake,
			BuildCommand:    fake,
			OutputDirectory: fake,
			StartCommand:    fake,
		}).Expect().Status(http.StatusOK)

	apps := expect.GET("/project/" + projectID + "/apps").
		Expect().
		Status(http.StatusOK).
		JSON()

	apps.Array().Length().Equal(1)

	appID := apps.Array().Element(0).Object().Value("id").String().Raw()

	expect.GET(fmt.Sprintf("/project/%s", projectID)).
		Expect().Status(http.StatusOK).JSON().
		Object().
		ContainsKey("name").ValueEqual("name", projectName)

	expect.GET(fmt.Sprintf("/app/%s", appID)).
		Expect().Status(http.StatusOK).JSON().
		Object().
		ContainsKey("projectID").ValueEqual("projectID", projectID).
		ContainsKey("name").ValueEqual("name", appName)

	expect.GET(fmt.Sprintf("/app/%s/deployments", appID)).
		Expect().Status(http.StatusOK).JSON().
		Null()

	expect.DELETE(fmt.Sprintf("/app/%s", appID)).
		Expect().Status(http.StatusOK)

	expect.DELETE(fmt.Sprintf("/project/%s", projectID)).
		Expect().Status(http.StatusOK)

	expect.GET(fmt.Sprintf("/project/%s", projectID)).
		Expect().Status(http.StatusOK).JSON().Null()

	expect.GET(fmt.Sprintf("/app/%s", appID)).
		Expect().Status(http.StatusOK).JSON().Null()
}
