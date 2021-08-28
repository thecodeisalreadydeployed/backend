package model

type DeploymentState string

const (
	DeploymentStateQueueing       DeploymentState = "DeploymentStateQueueing"
	DeploymentStateBuilding       DeploymentState = "DeploymentStateBuilding"
	DeploymentStateBuildSucceeded DeploymentState = "DeploymentStateBuildSucceeded"
	DeploymentStateReady          DeploymentState = "DeploymentStateReady"
	DeploymentStateError          DeploymentState = "DeploymentStateError"
)
