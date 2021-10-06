package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"github.com/thecodeisalreadydeployed/apiserver/dto"
	"github.com/thecodeisalreadydeployed/model"
	"io/ioutil"
	"net/http"
	"testing"
)

func setup(t *testing.T) *httpexpect.Expect {

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:3000",
		Reporter: httpexpect.NewAssertReporter(t),
	})

	return e
}

func httpRequest(path string, t *testing.T) []byte {
	resp, err := http.Get(fmt.Sprintf("http://localhost:3000%s", path))

	if err != nil {
		t.Error(err)
	}

	if resp.Status != "200 OK" {
		t.Error(fmt.Sprintf("Non-200 status code while requesting %s.", path))
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	return body
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

func TestFlow(t *testing.T) {
	projectName := "Test Project"
	appName := "Test App"
	fake := "Fake Data"

	expect := setup(t)

	expect.GET("/projects").Expect().Status(http.StatusOK).JSON().Null()
	expect.GET("/apps").Expect().Status(http.StatusOK).JSON().Null()

	expect.POST("/project").
		WithForm(dto.CreateProjectRequest{Name: projectName}).
		Expect().
		Status(http.StatusOK)

	var projects []model.Project
	bytes := httpRequest(fmt.Sprintf("/project/name/%s", projectName), t)
	err := json.Unmarshal(bytes, &projects)

	if err != nil {
		t.Error(err)
	}

	if len(projects) == 0 {
		t.Fatal("Test project was not created.")
	}
	project := projects[0]

	expect.POST("/app").
		WithForm(dto.CreateAppRequest{
			ProjectID:       project.ID,
			Name:            appName,
			RepositoryURL:   fake,
			BuildScript:     fake,
			InstallCommand:  fake,
			BuildCommand:    fake,
			OutputDirectory: fake,
			StartCommand:    fake,
		}).Expect().Status(http.StatusOK)

	var apps []model.App
	bytes = httpRequest(fmt.Sprintf("/app/name/%s", appName), t)
	err = json.Unmarshal(bytes, &apps)
	if err != nil {
		t.Error(err)
	}

	if len(apps) == 0 {
		t.Fatal("Test app was not created.")
	}
	app := apps[0]

	expect.GET(fmt.Sprintf("/project/%s", project.ID)).
		Expect().Status(http.StatusOK).JSON().Object().ContainsMap(project)

	expect.GET(fmt.Sprintf("/app/%s", app.ID)).
		Expect().Status(http.StatusOK).JSON().Object().ContainsMap(app)

	expect.GET(fmt.Sprintf("/project/%s/apps", project.ID)).
		Expect().Status(http.StatusOK).JSON().Array().Contains(app)

	expect.GET(fmt.Sprintf("/app/%s/deployments", app.ID)).
		Expect().Status(http.StatusOK).JSON().Null()

	expect.GET(fmt.Sprintf("/project/name/%s", projectName)).
		Expect().Status(http.StatusOK).JSON().Array().ContainsOnly(project)

	expect.GET(fmt.Sprintf("/app/name/%s", appName)).
		Expect().Status(http.StatusOK).JSON().Array().ContainsOnly(app)

	expect.DELETE(fmt.Sprintf("/app/%s", app.ID)).
		Expect().Status(http.StatusOK)

	expect.DELETE(fmt.Sprintf("/project/%s", project.ID)).
		Expect().Status(http.StatusOK)

	expect.GET(fmt.Sprintf("/project/name/%s", projectName)).
		Expect().Status(http.StatusOK).JSON().Null()

	expect.GET(fmt.Sprintf("/app/name/%s", appName)).
		Expect().Status(http.StatusOK).JSON().Null()
}
