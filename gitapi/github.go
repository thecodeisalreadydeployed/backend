package gitapi

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

func ListBranches(url string) []string {
	name, repo := getNameAndRepo(url)
	urlapi := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches", name, repo)
	res, err := http.Get(urlapi)
	if err != nil {
		zap.L().Error(err.Error())
	}
	if res == nil {
		return nil
	}

	var body []map[string]interface{}
	bytes, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		zap.L().Error(err.Error())
	}

	var output []string
	for _, branch := range body {
		output = append(output, cast.ToString(branch["name"]))
	}

	if res.Body != nil {
		err = res.Body.Close()
		if err != nil {
			zap.L().Error(err.Error())
		}
	}

	return output
}

// Returns name and repository, in order.
// For example, inputting "https://github.com/octocat/Hello-World"
// would return ("octocat", "Hello-World")
func getNameAndRepo(url string) (string, string) {
	urlslice := strings.Split(url, "/")
	return urlslice[len(urlslice)-2], urlslice[len(urlslice)-1]
}
