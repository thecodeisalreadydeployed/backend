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
	apiURL := config.ArgoCDServerHost() + "/api/v1/application?name=" + "codedeploy" + "&refresh=true"
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
