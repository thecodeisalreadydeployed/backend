package model

type DeploymentState string

const (
	DeploymentStateQueueing       DeploymentState = "DeploymentStateQueueing"
	DeploymentStateBuilding       DeploymentState = "DeploymentStateBuilding"
	DeploymentStateBuildSucceeded DeploymentState = "DeploymentStateBuildSucceeded"
	DeploymentStateCommited       DeploymentState = "DeploymentStateCommited"
	DeploymentStateReady          DeploymentState = "DeploymentStateReady"
	DeploymentStateError          DeploymentState = "DeploymentStateError"
)
