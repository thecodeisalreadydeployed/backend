package config

import containerregistry "github.com/thecodeisalreadydeployed/containerregistry/types"

type GitCredential struct {
	Repository string
	SSHKey     string
}

type BackendConfig struct {
	Registries     []containerregistry.ContainerRegistryConfiguration
	GitCredentials []GitCredential
}
