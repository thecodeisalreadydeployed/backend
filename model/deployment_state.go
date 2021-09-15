package model

type DeploymentState string

const (
	DeploymentStateQueueing       DeploymentState = "DeploymentStateQueueing"
	DeploymentStateBuilding       DeploymentState = "DeploymentStateBuilding"
	DeploymentStateBuildSucceeded DeploymentState = "DeploymentStateBuildSucceeded"
	DeploymentStateCommitted      DeploymentState = "DeploymentStateCommitted"
	DeploymentStateReady          DeploymentState = "DeploymentStateReady"
	DeploymentStateError          DeploymentState = "DeploymentStateError"
)
