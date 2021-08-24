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
