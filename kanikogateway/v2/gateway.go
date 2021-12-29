package kanikogateway

import (
	"github.com/thecodeisalreadydeployed/containerregistry"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/kubernetesinteractor/v2"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KanikoGateway interface {
	BuildContainerImage() (string, error)
}

type kanikoGateway struct {
	deploymentID       string
	repositoryURL      string
	branch             string
	buildConfiguration model.BuildConfiguration
	kubernetes         *kubernetesinteractor.KubernetesInteractor
	registry           *containerregistry.ContainerRegistry
}

const imageTag = "latest"

func NewKanikoGateway(deploymentID string, repositoryURL string, branch string, buildConfiguration model.BuildConfiguration) (KanikoGateway, error) {
	it, err := kubernetesinteractor.NewKubernetesInteractor()

	if err != nil {
		zap.L().Error("could not initialize KubernetesInteractor", zap.String("deploymentID", deploymentID))
		return nil, errutil.ErrFailedPrecondition
	}

	return kanikoGateway{deploymentID: deploymentID, repositoryURL: repositoryURL, branch: branch, buildConfiguration: buildConfiguration, kubernetes: &it, registry: nil}, nil
}

func (kanikoGateway) BuildContainerImage() (string, error) {
	return "", errutil.ErrNotImplemented
}

func (it kanikoGateway) kanikoPod() apiv1.Pod {
	workspace := apiv1.VolumeMount{
		Name:      "workspace",
		MountPath: "/workspace",
	}

	podLabel := map[string]string{
		"thecodeisalreadydeployed.github/deployment-id": it.deploymentID,
		"thecodeisalreadydeployed.github/component":     "KANIKO",
	}

	_ = it.buildConfiguration.BuildScript

	pod := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   it.deploymentID,
			Labels: podLabel,
		},
		Spec: apiv1.PodSpec{
			RestartPolicy: apiv1.RestartPolicyNever,
			Volumes: []apiv1.Volume{
				{
					Name: workspace.Name,
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
			},
			InitContainers: []apiv1.Container{
				{
					Name:         "imagebuilder-workspace",
					Image:        "ghcr.io/thecodeisalreadydeployed/imagebuilder-workspace:" + imageTag,
					VolumeMounts: []apiv1.VolumeMount{workspace},
				},
			},
		},
	}

	return pod
}
