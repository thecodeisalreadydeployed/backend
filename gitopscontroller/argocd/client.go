//go:generate mockgen -destination mock/client.go . ArgoCDClient

package argocd

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/thecodeisalreadydeployed/config"
)

type ArgoCDClient interface {
	Refresh() error
	Sync() error
}

type argoCDClient struct {
	httpClient *http.Client
}

func NewArgoCDClient() ArgoCDClient {
	var transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var httpClient = &http.Client{
		Timeout:   2 * time.Second,
		Transport: transport,
	}

	return &argoCDClient{httpClient: httpClient}
}

func (client *argoCDClient) Refresh() error {
	apiURL := config.ArgoCDServerHost() + "/api/v1/applications?name=" + "codedeploy" + "&refresh=true"
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
	apiURL := config.ArgoCDServerHost() + "/api/v1/applications/" + "codedeploy" + "/sync"
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
