package argocd

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/thecodeisalreadydeployed/config"
)

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var httpClient = &http.Client{
	Timeout:   2 * time.Second,
	Transport: transport,
}

func Refresh() error {
	apiURL := config.ArgoCDServerHost() + "/api/v1/application?name=" + "codedeploy" + "&refresh=true"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
