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
	Deploy() (string, error)
}

type kanikoGateway struct {
	projectID          string
	appID              string
	deploymentID       string
	repositoryURL      string
	branch             string
	buildConfiguration model.BuildConfiguration
	kubernetes         kubernetesinteractor.KubernetesInteractor
	registry           *containerregistry.ContainerRegistry
	logger             *zap.SugaredLogger
}

const imageTag = "latest"

func NewKanikoGateway(
	projectID string,
	appID string,
	deploymentID string,
	repositoryURL string,
	branch string,
	buildConfiguration model.BuildConfiguration,
	containerRegistry containerregistry.ContainerRegistry,
) (KanikoGateway, error) {
	logger := zap.L().Sugar().With("deploymentID", deploymentID)
	it, err := kubernetesinteractor.NewKubernetesInteractor()

	if err != nil {
		logger.Error("failed to initialize KubernetesInteractor")
		return nil, errutil.ErrFailedPrecondition
	}

	return kanikoGateway{
		projectID:          projectID,
		appID:              appID,
		deploymentID:       deploymentID,
		repositoryURL:      repositoryURL,
		branch:             branch,
		buildConfiguration: buildConfiguration,
		kubernetes:         it,
		registry:           &containerRegistry,
		logger:             logger,
	}, nil
}

func (it kanikoGateway) Deploy() (string, error) {
	workspace := apiv1.VolumeMount{
		Name:      "workspace",
		MountPath: "/workspace",
	}

	objectLabel := map[string]string{
		"deployment.api.deploys.dev/id":        it.deploymentID,
		"deployment.api.deploys.dev/component": "imagebuilder",
	}

	it.logger.Info(objectLabel)

	buildScript := it.buildConfiguration.BuildScript

	configMap := apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "imagebuilder-" + it.deploymentID,
			Labels: objectLabel,
		},
		Data: map[string]string{
			"Dockerfile": buildScript,
		},
	}

	it.logger.Info(configMap)

	pod := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "imagebuilder-" + it.deploymentID,
			Labels: objectLabel,
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
					Env: []apiv1.EnvVar{
						{Name: "CODEDEPLOY_DEPLOYMENT_ID", Value: it.deploymentID},
						{Name: "CODEDEPLOY_GIT_REPOSITORY", Value: it.repositoryURL},
						{Name: "CODEDEPLOY_GIT_REFERENCE", Value: it.branch},
					},
				},
			},
			Containers: []apiv1.Container{
				{
					Name:         "imagebuilder",
					Image:        "ghcr.io/thecodeisalreadydeployed/imagebuilder:" + imageTag,
					VolumeMounts: []apiv1.VolumeMount{workspace},
					Env: []apiv1.EnvVar{
						{Name: "CODEDEPLOY_DEPLOYMENT_ID", Value: it.deploymentID},
						{Name: "CODEDEPLOY_KANIKO_LOG_VERBOSITY", Value: "info"},
						{Name: "CODEDEPLOY_KANIKO_CONTEXT", Value: "/workspace/" + it.deploymentID},
					},
				},
			},
		},
	}

	it.logger.Info(pod)

	_, err := it.kubernetes.CreatePod(pod, it.projectID)
	if err != nil {
		it.logger.Error("failed to create imagebuilder pod")
	}

	return it.deploymentID, nil
}
