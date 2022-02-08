package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/gitapi/provider"
	"go.uber.org/zap"
)

type gitHubAPI struct {
	logger *zap.Logger
	owner  string
	repo   string
}

func NewGitHubAPI(logger *zap.Logger, owner string, repo string) provider.GitProvider {
	return &gitHubAPI{logger: logger, owner: owner, repo: repo}
}

// List branches in strings given a GitHub utl string.
func (gh *gitHubAPI) GetBranches() ([]string, error) {
	urlapi := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches", gh.owner, gh.repo)
	res, err := http.Get(urlapi)
	defer closeHTTP(res)
	if err != nil {
		gh.logger.Error(err.Error())
		return nil, errutil.ErrUnavailable
	}

	var body []Branch
	err = getJSON(res, &body)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrUnknown
	}

	var output []string
	for _, branch := range body {
		output = append(output, branch.Name)
	}

	return output, nil
}

// List all file names in strings given a GitHub url string and branch name.
func (gh *gitHubAPI) GetFiles(branch string) ([]string, error) {
	urlapi := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1", gh.owner, gh.repo, branch)
	res, err := http.Get(urlapi)
	defer closeHTTP(res)
	if err != nil {
		gh.logger.Error(err.Error())
		return nil, errutil.ErrUnavailable
	}

	var body File
	err = getJSON(res, &body)
	if err != nil {
		gh.logger.Error(err.Error())
		return nil, errutil.ErrUnknown
	}

	var output []string
	for _, file := range body.Tree {
		if cast.ToString(file.Type) == "blob" {
			output = append(output, file.Path)
		}
	}

	return output, nil
}

// Get raw file given GitHub url string, branch, and file path.
func (gh *gitHubAPI) GetRaw(branch string, path string) (string, error) {
	urlapi := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", gh.owner, gh.repo, branch, path)
	res, err := http.Get(urlapi)
	defer closeHTTP(res)
	if err != nil {
		gh.logger.Error(err.Error())
		return "", errutil.ErrUnavailable
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		gh.logger.Error(err.Error())
		return "", errutil.ErrUnknown
	}
	return string(bytes), nil
}

// Gets JSON from HTTP response.
func getJSON(res *http.Response, output interface{}) error {
	return json.NewDecoder(res.Body).Decode(&output)
}

// Close HTTP connection.
func closeHTTP(res *http.Response) {
	err := res.Body.Close()
	if err != nil {
		zap.L().Error(err.Error())
	}
}
