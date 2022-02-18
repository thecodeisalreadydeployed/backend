package kanikogateway

import (
	"github.com/google/uuid"
	"github.com/thecodeisalreadydeployed/clusterbackend"
	containerregistry "github.com/thecodeisalreadydeployed/containerregistry/types"
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
	clusterBackend     clusterbackend.ClusterBackend
	registry           containerregistry.ContainerRegistry
	logger             *zap.Logger
}

const codeDeployInternalNamespace = "codedeploy-internal"
const imageTag = "latest"

func NewKanikoGateway(
	logger *zap.Logger,
	clusterBackend clusterbackend.ClusterBackend,
	projectID string,
	appID string,
	deploymentID string,
	repositoryURL string,
	branch string,
	buildConfiguration model.BuildConfiguration,
	containerRegistry containerregistry.ContainerRegistry,
) (KanikoGateway, error) {
	return kanikoGateway{
		projectID:          projectID,
		appID:              appID,
		deploymentID:       deploymentID,
		repositoryURL:      repositoryURL,
		branch:             branch,
		buildConfiguration: buildConfiguration,
		clusterBackend:     clusterBackend,
		registry:           containerRegistry,
		logger:             logger.With(zap.String("deploymentID", deploymentID)),
	}, nil
}

func (gateway kanikoGateway) Deploy() (string, error) {
	workspace := apiv1.VolumeMount{
		Name:      "workspace",
		MountPath: "/workspace",
	}

	UID := uuid.NewString()

	objectLabel := map[string]string{
		"beta.deploys.dev/uid":           UID,
		"beta.deploys.dev/deployment-id": gateway.deploymentID,
		"beta.deploys.dev/component":     "imagebuilder",
	}

	configMap := apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "builder-" + UID,
			Labels: objectLabel,
		},
		Data: map[string]string{
			"Dockerfile": gateway.buildConfiguration.BuildScript,
		},
	}

	kanikoDestination := ""
	kubernetesServiceAccountName := ""
	containerRegistry := gateway.registry
	kanikoDestination = containerRegistry.RegistryFormat(gateway.appID, gateway.deploymentID)

	if containerRegistry.AuthenticationMethod() == containerregistry.KubernetesServiceAccount {
		kubernetesServiceAccountName = containerRegistry.Secret()
	}

	pod := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "builder-" + UID,
			Labels: objectLabel,
		},
		Spec: apiv1.PodSpec{
			ServiceAccountName: kubernetesServiceAccountName,
			RestartPolicy:      apiv1.RestartPolicyNever,
			Volumes: []apiv1.Volume{
				{
					Name: workspace.Name,
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
				{
					Name: "dockerfile",
					VolumeSource: apiv1.VolumeSource{
						ConfigMap: &apiv1.ConfigMapVolumeSource{
							LocalObjectReference: apiv1.LocalObjectReference{
								Name: "builder-" + UID,
							},
						},
					},
				},
			},
			InitContainers: []apiv1.Container{
				{
					Name:    "workload-identity-init-container",
					Image:   "gcr.io/google.com/cloudsdktool/cloud-sdk:326.0.0-alpine",
					Command: []string{"/bin/bash", "-c", "curl -s -H 'Metadata-Flavor: Google' 'http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token' --retry 30 --retry-connrefused --retry-max-time 30 > /dev/null || exit 1"},
				},
				{
					Name:         "imagebuilder-workspace",
					Image:        "ghcr.io/thecodeisalreadydeployed/imagebuilder-workspace:" + imageTag,
					VolumeMounts: []apiv1.VolumeMount{workspace},
					Env: []apiv1.EnvVar{
						// TODO(trif0lium): use environment variable
						{Name: "CODEDEPLOY_API_URL", Value: "http://codedeploy.default.svc.cluster.local:3000"},
						{Name: "CODEDEPLOY_DEPLOYMENT_ID", Value: gateway.deploymentID},
						{Name: "CODEDEPLOY_GIT_REPOSITORY", Value: gateway.repositoryURL},
						{Name: "CODEDEPLOY_GIT_REFERENCE", Value: gateway.branch},
					},
				},
			},
			Containers: []apiv1.Container{
				{
					Name:  "imagebuilder",
					Image: "ghcr.io/thecodeisalreadydeployed/imagebuilder:" + imageTag,
					VolumeMounts: []apiv1.VolumeMount{
						workspace,
						{
							Name:      "dockerfile",
							MountPath: "/kaniko/deploys-dev/Dockerfile",
							SubPath:   "Dockerfile",
						},
					},
					Env: []apiv1.EnvVar{
						// TODO(trif0lium): use environment variable
						{Name: "CODEDEPLOY_API_URL", Value: "http://codedeploy.default.svc.cluster.local:3000"},
						{Name: "CODEDEPLOY_DEPLOYMENT_ID", Value: gateway.deploymentID},
						{Name: "CODEDEPLOY_KANIKO_LOG_VERBOSITY", Value: "info"},
						{Name: "CODEDEPLOY_KANIKO_CONTEXT", Value: "/workspace/" + gateway.deploymentID},
						{Name: "CODEDEPLOY_KANIKO_DESTINATION", Value: kanikoDestination},
					},
				},
			},
		},
	}

	_, err := gateway.clusterBackend.CreateConfigMap(configMap, codeDeployInternalNamespace)
	if err != nil {
		return "", err
	}

	_, err = gateway.clusterBackend.CreatePod(pod, codeDeployInternalNamespace)
	if err != nil {
		return "", err
	}

	return gateway.deploymentID, nil
}
