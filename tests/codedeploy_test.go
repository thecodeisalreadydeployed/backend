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

func httpRequest(path string, obj interface{}, t *testing.T) interface{} {
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

	err = json.Unmarshal(body, obj)

	if err != nil {
		t.Error(err)
	}

	return obj
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

	expect.POST("/project").
		WithForm(dto.CreateProjectRequest{Name: projectName}).
		Expect().
		Status(http.StatusOK)

	prj := httpRequest(fmt.Sprintf("/project/name/%s", projectName), model.Project{}, t).(model.Project)

	expect.POST("/app").
		WithForm(dto.CreateAppRequest{
			ProjectID:       prj.ID,
			Name:            appName,
			RepositoryURL:   fake,
			BuildScript:     fake,
			InstallCommand:  fake,
			BuildCommand:    fake,
			OutputDirectory: fake,
			StartCommand:    fake,
		}).Expect().
		Status(http.StatusOK)

	app := httpRequest(fmt.Sprintf())

	expect.GET(fmt.Sprintf("/project/%s", prj.ID)).
		Expect().Status(http.StatusOK).JSON().Object().ContainsMap(prj)
}
