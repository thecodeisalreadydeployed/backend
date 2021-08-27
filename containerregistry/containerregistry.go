package containerregistry

type ContainerRegistryType string

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

type ContainerRegistry interface {
	RegistryFormat(repository string, tag string) (string, error)
	Type() ContainerRegistryType
	Secret() string
}
