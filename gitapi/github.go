package gitapi

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/errutil"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// List branches in strings given a GitHub utl string.
func GetBranches(url string) ([]string, error) {
	name, repo := getNameAndRepo(url)
	urlapi := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches", name, repo)
	res, err := http.Get(urlapi)
	defer closeHTTP(res)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrUnavailable
	}

	body, err := getJSONArray(res)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrUnknown
	}

	var output []string
	for _, branch := range body {
		output = append(output, cast.ToString(branch["name"]))
	}

	return output, nil
}

// List all file names in strings given a GitHub url string and branch name.
func GetFiles(url string, branch string) ([]string, error) {
	name, repo := getNameAndRepo(url)
	urlapi := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1", name, repo, branch)
	res, err := http.Get(urlapi)
	defer closeHTTP(res)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrUnavailable
	}

	body, err := getJSONArray(res)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrUnknown
	}

	var output []string
	for _, file := range body {
		if cast.ToString(file["type"]) == "file" {
			output = append(output, cast.ToString(file["path"]))
		}
	}

	return output, nil
}

// Get raw file given GitHub url string, branch, and file path.
func GetRaw(url string, branch string, path string) (string, error) {
	name, repo := getNameAndRepo(url)
	urlapi := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", name, repo, branch, path)
	res, err := http.Get(urlapi)
	defer closeHTTP(res)
	if err != nil {
		zap.L().Error(err.Error())
		return "", errutil.ErrUnavailable
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error(err.Error())
		return "", errutil.ErrUnknown
	}
	return string(bytes), nil
}

// Returns name and repository, in order. Must have HTTPS prefix.
// For example, inputting "https://github.com/octocat/Hello-World"
// would return ("octocat", "Hello-World")
func getNameAndRepo(url string) (string, string) {
	urlslice := strings.Split(url, "/")
	return urlslice[3], urlslice[4]
}

// Gets JSON from HTTP response.
func getJSONArray(res *http.Response) ([]map[string]interface{}, error) {
	var body []map[string]interface{}
	bytes, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		return nil, errutil.ErrUnknown
	}
	return body, nil
}

// Close HTTP connection.
func closeHTTP(res *http.Response) {
	err := res.Body.Close()
	if err != nil {
		zap.L().Error(err.Error())
	}
}
