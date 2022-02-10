//go:generate mockgen -destination mock/client.go . ArgoCDClient

package argocd

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/thecodeisalreadydeployed/config"
	"go.uber.org/zap"
)

type ArgoCDClient interface {
	Refresh() error
	Sync() error
}

type argoCDClient struct {
	httpClient *http.Client

	logger   *zap.Logger
	appName  string
	repoPath string
}

func NewArgoCDClient(logger *zap.Logger, appName string, repoPath string) ArgoCDClient {
	var transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var httpClient = &http.Client{
		Timeout:   2 * time.Second,
		Transport: transport,
	}

	return &argoCDClient{httpClient: httpClient, logger: logger, appName: appName, repoPath: repoPath}
}

func (client *argoCDClient) Refresh() error {
	apiURL := config.ArgoCDServerHost() + "/api/v1/applications?name=" + client.appName + "&refresh=true"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *argoCDClient) Sync() error {
	apiURL := config.ArgoCDServerHost() + "/api/v1/applications/" + client.appName + "/sync"
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
