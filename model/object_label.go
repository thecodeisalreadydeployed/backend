package model

func DeploymentLabel(deploymentID string) map[string]string {
	return map[string]string{
		"codedeploy/deployment-id": deploymentID,
	}
}

func PodLabel(deploymentID string) map[string]string {
	return map[string]string{
		"codedeploy/deployment-id": deploymentID,
	}
}
