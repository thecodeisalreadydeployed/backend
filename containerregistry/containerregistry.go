package containerregistry

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
	KubernetesServiceAccount() string
}

type ContainerRegistryConfiguration struct {
	Type                 ContainerRegistryType `json:"type"`
	AuthenticationMethod AuthenticationMethod  `json:"authenticationMethod"`
	Repository           string                `json:"repository"`
	Secret               string                `json:"secret"`
}
