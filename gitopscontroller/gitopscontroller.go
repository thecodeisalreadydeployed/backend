package gitopscontroller

type GitOpsController interface {
	SetupUserspace()
	SetupProject(projectID string) error
	SetupApp(projectID string, appID string) error
	UpdateContainerImage(appID string, deploymentID string) error
}
