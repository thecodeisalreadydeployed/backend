//go:generate mockgen -destination mock/client.go . ArgoCDClient

package argocd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/util"
	"go.uber.org/zap"
)

type ArgoCDClient interface {
	CreateApp() error
	Refresh() error
	Sync() error
}

type argoCDClient struct {
	httpClient *http.Client

	logger   *zap.Logger
	appName  string
	repoPath string

	isInitialized bool
}

var HTTPTransport http.RoundTripper = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func NewArgoCDClient(logger *zap.Logger, appName string, repoPath string) ArgoCDClient {
	var httpClient = &http.Client{
		Transport: HTTPTransport,
	}

	isInitialized := true

	if util.IsDevEnvironment() {
		isInitialized = false
	}

	if util.IsDockerTestEnvironment() {
		isInitialized = false
	}

	return &argoCDClient{httpClient: httpClient, logger: logger.With(zap.String("appName", appName), zap.String("repoPath", repoPath)), appName: appName, repoPath: repoPath, isInitialized: isInitialized}
}

func (client *argoCDClient) CreateApp() error {
	if !client.isInitialized {
		return nil
	}

	apiURL := config.ArgoCDServerHost() + "/api/v1/applications"
	u, _ := url.Parse(config.GitServerHost())
	u.Path = client.appName
	repoURL := u.String()
	requestBody, _ := json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":       client.appName,
			"namespace":  "argocd",
			"finalizers": []string{"resources-finalizer.argocd.argoproj.io"},
		},
		"spec": map[string]interface{}{
			"project": "default",
			"source": map[string]string{
				"path":           ".",
				"repoURL":        repoURL,
				"targetRevision": "master",
			},
			"destination": map[string]string{
				"server":    "https://kubernetes.default.svc",
				"namespace": "default",
			},
			"syncPolicy": map[string]interface{}{
				"automated": map[string]bool{
					"prune":    true,
					"selfHeal": true,
				},
				"syncOptions": []string{"CreateNamespace=true"},
				"retry": map[string]interface{}{
					"limit": 5,
					"backoff": map[string]interface{}{
						"duration":    "5s",
						"factor":      2,
						"maxDuration": "3m",
					},
				},
			},
		},
	})

	requestID := uuid.NewString()
	client.logger.Info("creating Argo CD application", zap.String("requestBody", string(requestBody)))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		client.logger.Error(err.Error(), zap.String("requestID", requestID))
		return err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		client.logger.Error(err.Error(), zap.String("requestID", requestID))
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.logger.Error(err.Error(), zap.String("requestID", requestID))
		return err
	}

	client.logger.Info(string(responseBody), zap.String("requestID", requestID))

	return nil
}

func (client *argoCDClient) Refresh() error {
	if !client.isInitialized {
		return nil
	}

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
	if !client.isInitialized {
		return nil
	}

	requestID := uuid.NewString()
	client.logger.Info("syncing Argo CD application")

	apiURL := config.ArgoCDServerHost() + "/api/v1/applications/" + client.appName + "/sync"
	requestBody, _ := json.Marshal(map[string]interface{}{
		"dryRun":    false,
		"prune":     true,
		"resources": nil,
		"revision":  "master",
		"strategy": map[string]interface{}{
			"apply": map[string]bool{
				"force": true,
			},
		},
		"syncOptions": map[string]interface{}{
			"items": []string{"CreateNamespace=true"},
		},
	})

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		client.logger.Error(err.Error(), zap.String("requestID", requestID))
		return err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		client.logger.Error(err.Error(), zap.String("requestID", requestID))
		return err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.logger.Error(err.Error(), zap.String("requestID", requestID))
		return err
	}

	client.logger.Info(string(responseBody), zap.String("requestID", requestID))
	return nil
}
