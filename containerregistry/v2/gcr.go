package containerregistry

import (
	"context"
	"fmt"

	"github.com/thecodeisalreadydeployed/errutil"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
)

type gcrGateway struct {
	hostname          string
	projectID         string
	repository        string
	serviceAccountKey string
}

func NewGCRGateway(hostname string, projectID string, repository string, serviceAccountKey string) (ContainerRegistry, error) {
	c, err := google.DefaultClient(context.TODO(), cloudresourcemanager.CloudPlatformScope)
	if err != nil {
		zap.L().Error("cannot initialize google.DefaultClient")
		return nil, errutil.ErrFailedPrecondition
	}

	cloudResourceManagerService, err := cloudresourcemanager.New(c)
	if err != nil {
		zap.L().Error("cannot initialize Cloud Resource Manager")
		return nil, errutil.ErrFailedPrecondition
	}

	resource := "artifactregistry"

	permissions := []string{
		"artifactregistry.repositories.uploadArtifacts",
		"artifactregistry.tags.create",
		"artifactregistry.tags.update",
	}

	request := &cloudresourcemanager.TestIamPermissionsRequest{Permissions: permissions}

	resp, err := cloudResourceManagerService.Projects.TestIamPermissions(resource, request).Context(context.TODO()).Do()
	if err != nil {
		zap.L().Error("cannot call TestIamPermissions")
		return nil, errutil.ErrUnavailable
	}

	if len(resp.Permissions) == len(permissions) {
		return gcrGateway{hostname: hostname, projectID: projectID, repository: repository, serviceAccountKey: serviceAccountKey}, nil
	}

	return nil, errutil.ErrUnavailable
}

func (gateway gcrGateway) ImageName(name string, tag string) string {
	return fmt.Sprintf("%s/%s/%s/%s:%s", gateway.hostname, gateway.projectID, gateway.repository, name, tag)
}

func (gateway gcrGateway) Type() ContainerRegistryType {
	return GCR
}
