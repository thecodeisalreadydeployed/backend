package config

import "github.com/thecodeisalreadydeployed/containerregistry"

type GitCredential struct {
	Repository string
	SSHKey     string
}

type BackendConfig struct {
	Registries     []containerregistry.ContainerRegistryConfiguration
	GitCredentials []GitCredential
}
