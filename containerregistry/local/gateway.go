package gcr

import (
	"fmt"

	"github.com/thecodeisalreadydeployed/containerregistry/types"
)

func NewLocalRegistryGateway(hostname string, port int, authenticationMethod types.AuthenticationMethod, secret string) types.ContainerRegistry {
	return &localRegistryGateway{hostname: hostname, port: port, authenticationMethod: authenticationMethod, secret: secret}
}

type localRegistryGateway struct {
	hostname             string
	port                 int
	authenticationMethod types.AuthenticationMethod
	secret               string
}

func (registry *localRegistryGateway) RegistryFormat(image string, tag string) string {
	return fmt.Sprintf("%s:%d/%s:%s", registry.hostname, registry.port, image, tag)
}

func (registry *localRegistryGateway) Type() types.ContainerRegistryType {
	return types.LOCAL
}

func (registry *localRegistryGateway) Secret() string {
	return registry.secret
}

func (registry *localRegistryGateway) AuthenticationMethod() types.AuthenticationMethod {
	return registry.authenticationMethod
}
