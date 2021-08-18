package model

type DeploymentState string

const (
	DeploymentStateQueueing DeploymentState = "DeploymentStateQueueing"
	DeploymentStateBuilding DeploymentState = "DeploymentStateBuilding"
	DeploymentStateReady    DeploymentState = "DeploymentStateReady"
	DeploymentStateError    DeploymentState = "DeploymentStateError"
)
