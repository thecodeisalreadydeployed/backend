package gcr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thecodeisalreadydeployed/containerregistry"
)

func NewGCRGateway(hostname string, projectID string) containerregistry.ContainerRegistry {
	return &gcrGateway{hostname: hostname, projectID: projectID}
}

type gcrGateway struct {
	hostname  string
	projectID string
}

func (g *gcrGateway) RegistryFormat(repository string, tag string) (string, error) {
	if !strings.Contains(g.hostname, "gcr.io") {
		return "", errors.New("Invalid hostname for Google Container Registry.")
	}

	return fmt.Sprintf("%s/%s/%s:%s", g.hostname, g.projectID, repository, tag), nil
}
