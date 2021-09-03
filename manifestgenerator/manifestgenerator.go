package manifestgenerator

import "go.uber.org/zap"

func GenerateManifests() ([]string, error) {
	deploymentYAML, err := GenerateDeploymentYAML(&GenerateDeploymentOptions{
		Name:           "",
		Namespace:      "",
		Labels:         nil,
		ContainerImage: "",
	})

	if err != nil {
		zap.L().Error("Failed to generate deployment YAML.")
		return nil, err
	}

	serviceYAML, err := GenerateServiceYAML(&GenerateServiceOptions{
		Name:           "",
		Namespace:      "",
		Labels:         nil,
		ContainerImage: "",
	})

	if err != nil {
		zap.L().Error("Failed to generate service YAML.")
		return nil, err
	}

	virtualServerYAML, err := GenerateVirtualServerYAML(&GenerateVirtualServerOptions{
		Labels:    nil,
		ProjectID: "",
		AppID:     "",
	})

	if err != nil {
		zap.L().Error("Failed to generate service YAML.")
		return nil, err
	}

	return []string{deploymentYAML, serviceYAML, virtualServerYAML}, nil
}
