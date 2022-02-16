package statusapi

import (
	"crypto/tls"
	"net/http"

	"go.uber.org/zap"
)

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
	return "", nil
}
