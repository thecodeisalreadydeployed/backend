package statusapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/thecodeisalreadydeployed/datastore"
	"go.uber.org/zap"
)

type Metadata struct {
	NodeName     string
	PodName      string
	PodNamespace string
	ProjectID    string
	AppID        string
	DeploymentID string
}

type StatusAPIBackend interface {
	GetActiveDeploymentID(appID string) (string, error)
}

type statusAPIBackend struct {
	logger     *zap.Logger
	httpClient *http.Client
}

var HTTPTransport http.RoundTripper = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func NewStatusAPIBackend(logger *zap.Logger) StatusAPIBackend {
	return &statusAPIBackend{logger: logger, httpClient: &http.Client{Transport: HTTPTransport}}
}

func (backend *statusAPIBackend) GetActiveDeploymentID(appID string) (string, error) {
	requestID := uuid.NewString()
	logger := backend.logger.With(zap.String("appID", appID), zap.String("requestID", requestID))
	_ = logger

	app, err := datastore.GetAppByID(datastore.GetDB(), appID)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	metadataAPI := fmt.Sprintf("http://%s.%s.svc.cluster.local:5000", app.ID, app.ProjectID)
	req, err := http.NewRequest("GET", metadataAPI, nil)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	resp, err := backend.httpClient.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	var metadata Metadata
	err = json.Unmarshal(responseBody, &metadata)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	return metadata.DeploymentID, nil
}
