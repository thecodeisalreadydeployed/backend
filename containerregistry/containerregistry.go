package containerregistry

import (
	"github.com/thecodeisalreadydeployed/containerregistry/gcr"
	"go.uber.org/zap"
)

type ContainerRegistryType string
type AuthenticationMethod string

const (
	// Local
	LOCAL ContainerRegistryType = "LOCAL"

	// Docker Hub
	DH ContainerRegistryType = "DH"

	// GitHub Container Registry
	GHCR ContainerRegistryType = "GHCR"

	// Google Container Registry
	GCR ContainerRegistryType = "GCR"

	// Amazon Elastic Container Registry
	ECR ContainerRegistryType = "ECR"
)

const (
	KubernetesServiceAccount AuthenticationMethod = "KubernetesServiceAccount"
	Secret                   AuthenticationMethod = "Secret"
)

type ContainerRegistry interface {
	RegistryFormat(repository string, tag string) string
	Type() ContainerRegistryType
	AuthenticationMethod() AuthenticationMethod
	Secret() string
}

type ContainerRegistryConfiguration struct {
	Type                 ContainerRegistryType `json:"type"`
	AuthenticationMethod AuthenticationMethod  `json:"authenticationMethod"`
	Repository           string                `json:"repository"`
	Secret               string                `json:"secret"`
	Metadata             map[string]string     `json:"metadata"`
}

func NewContainerRegistry(config ContainerRegistryConfiguration) ContainerRegistry {
	if config.Type == GCR {
		if len(config.Metadata["GOOGLE_CLOUD_PROJECT"]) == 0 {
			zap.L().Fatal("missing required metadata GOOGLE_CLOUD_PROJECT", zap.String("type", string(config.Type)))
		}

		containerRegistry := gcr.NewGCRGateway(
			"asia-southeast1-docker.pkg.dev",
			config.Metadata["GOOGLE_CLOUD_PROJECT"],
			config.AuthenticationMethod,
			config.Secret,
		)

		return containerRegistry
	}

	zap.L().Fatal("unsupported container registry", zap.String("type", string(config.Type)))
	panic("unsupported container registry")
}
