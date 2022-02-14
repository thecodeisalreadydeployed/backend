package gcr

import (
	"fmt"

	"github.com/thecodeisalreadydeployed/containerregistry/types"
)

func NewGCRGateway(hostname string, projectID string, authenticationMethod types.AuthenticationMethod, secret string) types.ContainerRegistry {
	return &gcrGateway{hostname: hostname, projectID: projectID, authenticationMethod: authenticationMethod, secret: secret}
}

type gcrGateway struct {
	hostname             string
	projectID            string
	authenticationMethod types.AuthenticationMethod
	secret               string
}

func (gcr *gcrGateway) RegistryFormat(repository string, tag string) string {
	return fmt.Sprintf("%s/%s/codedeploy/%s:%s", gcr.hostname, gcr.projectID, repository, tag)
}

func (gcr *gcrGateway) Type() types.ContainerRegistryType {
	return types.GCR
}

func (gcr *gcrGateway) Secret() string {
	return gcr.secret
}

func (gcr *gcrGateway) AuthenticationMethod() types.AuthenticationMethod {
	return gcr.authenticationMethod
}
