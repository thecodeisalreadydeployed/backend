package containerregistry

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/containerregistry/gcr"
	"github.com/thecodeisalreadydeployed/containerregistry/local"
	"github.com/thecodeisalreadydeployed/containerregistry/types"
	"go.uber.org/zap"
)

func NewContainerRegistry(config types.ContainerRegistryConfiguration) types.ContainerRegistry {
	if config.Type == types.GCR {
		if len(config.Metadata["GOOGLE_CLOUD_PROJECT"]) == 0 {
			zap.L().Fatal("missing required metadata GOOGLE_CLOUD_PROJECT", zap.String("type", string(config.Type)))
		}

		if len(config.Metadata["GOOGLE_CLOUD_REGION"]) == 0 {
			zap.L().Fatal("missing required metadata GOOGLE_CLOUD_REGION", zap.String("type", string(config.Type)))
		}

		containerRegistry := gcr.NewGCRGateway(
			fmt.Sprintf("%s-docker.pkg.dev", config.Metadata["GOOGLE_CLOUD_REGION"]),
			config.Metadata["GOOGLE_CLOUD_PROJECT"],
			config.Repository,
			config.AuthenticationMethod,
			config.Secret,
		)

		return containerRegistry
	}

	if config.Type == types.LOCAL {
		if len(config.Metadata["HOSTNAME"]) == 0 {
			zap.L().Fatal("missing required metadata HOSTNAME", zap.String("type", string(config.Type)))
		}

		if len(config.Metadata["PORT"]) == 0 {
			zap.L().Fatal("missing required metadata PORT", zap.String("type", string(config.Type)))
		}

		containerRegistry := local.NewLocalRegistryGateway(
			config.Metadata["HOSTNAME"],
			cast.ToInt(config.Metadata["PORT"]),
			config.AuthenticationMethod,
			config.Secret,
		)

		return containerRegistry
	}

	zap.L().Fatal("unsupported container registry", zap.String("type", string(config.Type)))
	panic("unsupported container registry")
}
