package containerregistry

import (
	"context"

	"github.com/thecodeisalreadydeployed/errutil"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
)

type gcrGateway struct {
}

func NewGCRGateway() (ContainerRegistry, error) {
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

	resource := ""
	request := &cloudresourcemanager.TestIamPermissionsRequest{}
	permissions := []string{}

	resp, err := cloudResourceManagerService.Projects.TestIamPermissions(resource, request).Context(context.TODO()).Do()
	if err != nil {
		zap.L().Error("cannot call TestIamPermissions")
		return nil, errutil.ErrUnavailable
	}

	if len(resp.Permissions) == len(permissions) {
		return gcrGateway{}, nil
	}

	return nil, errutil.ErrUnavailable
}

func (gateway gcrGateway) ImageURL(name string, tag string) string {
	return ""
}
