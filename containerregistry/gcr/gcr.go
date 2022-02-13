package gcr

import (
	"fmt"

	"github.com/thecodeisalreadydeployed/containerregistry"
)

func NewGCRGateway(hostname string, projectID string, authenticationMethod containerregistry.AuthenticationMethod, serviceAccountKey string, kubernetesServiceAccount string) containerregistry.ContainerRegistry {
	return &gcrGateway{hostname: hostname, projectID: projectID, authenticationMethod: authenticationMethod, serviceAccountKey: serviceAccountKey, kubernetesServiceAccount: kubernetesServiceAccount}
}

type gcrGateway struct {
	hostname                 string
	projectID                string
	authenticationMethod     containerregistry.AuthenticationMethod
	serviceAccountKey        string
	kubernetesServiceAccount string
}

func (gcr *gcrGateway) RegistryFormat(repository string, tag string) string {
	return fmt.Sprintf("%s/%s/%s:%s", gcr.hostname, gcr.projectID, repository, tag)
}

func (gcr *gcrGateway) Type() containerregistry.ContainerRegistryType {
	return containerregistry.GCR
}

func (gcr *gcrGateway) Secret() string {
	return gcr.serviceAccountKey
}

func (gcr *gcrGateway) AuthenticationMethod() containerregistry.AuthenticationMethod {
	return gcr.authenticationMethod
}

func (gcr *gcrGateway) KubernetesServiceAccount() string {
	return gcr.kubernetesServiceAccount
}
