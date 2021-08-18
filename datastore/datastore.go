package datastore

import "github.com/thecodeisalreadydeployed/model"

func GetProject(p *model.Project) *model.Project {
	return new(model.Project)
}

func GetApp(app *model.App) *model.App {
	return new(model.App)
}

func GetDeployment(dpl *model.Deployment) *model.Deployment {
	return new(model.Deployment)
}

func GetEvent(id string) string {
	return "Dummy event."
}
