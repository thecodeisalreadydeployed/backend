package manifestgenerator

import "go.uber.org/zap"

func GenerateManifests() (map[string]string, error) {
	deploymentYAML, err := GenerateDeploymentYAML(&GenerateDeploymentOptions{
		Name:           "DeploymentName",
		Namespace:      "DeploymentNamespace",
		Labels:         map[string]string{"DeploymentLabelKey": "DeploymentLabelValue"},
		ContainerImage: "DeploymentContainerImage",
	})

	if err != nil {
		zap.L().Error("Failed to generate deployment YAML.")
		return nil, err
	}

	serviceYAML, err := GenerateServiceYAML(&GenerateServiceOptions{
		Name:      "ServiceName",
		Namespace: "ServiceNamespace",
		Labels:    map[string]string{"ServiceLabelKey": "ServiceLabelValue"},
	})

	if err != nil {
		zap.L().Error("Failed to generate service YAML.")
		return nil, err
	}

	virtualServerYAML, err := GenerateVirtualServerYAML(&GenerateVirtualServerOptions{
		Labels:    map[string]string{"VirtualServerLabelKey": "VirtualServerLabelValue"},
		ProjectID: "ProjectID",
		AppID:     "AppID",
	})

	if err != nil {
		zap.L().Error("Failed to generate service YAML.")
		return nil, err
	}

	return map[string]string{
		"deploymentYAML":    deploymentYAML,
		"serviceYAML":       serviceYAML,
		"virtualServerYAML": virtualServerYAML,
	}, nil
}
