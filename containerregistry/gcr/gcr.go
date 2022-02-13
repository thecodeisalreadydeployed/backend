package gcr

import (
	"fmt"

	"github.com/thecodeisalreadydeployed/containerregistry"
)

func NewGCRGateway(hostname string, projectID string, serviceAccountKey string) containerregistry.ContainerRegistry {
	return &gcrGateway{hostname: hostname, projectID: projectID, serviceAccountKey: serviceAccountKey}
}

type gcrGateway struct {
	hostname          string
	projectID         string
	serviceAccountKey string
}

func (g *gcrGateway) RegistryFormat(repository string, tag string) string {
	return fmt.Sprintf("%s/%s/%s:%s", g.hostname, g.projectID, repository, tag)
}

func (g *gcrGateway) Type() containerregistry.ContainerRegistryType {
	return containerregistry.GCR
}

func (g *gcrGateway) Secret() string {
	return g.serviceAccountKey
}
