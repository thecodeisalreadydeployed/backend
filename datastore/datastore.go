package datastore

import "github.com/thecodeisalreadydeployed/model"

func GetProject(p *model.Project) *model.Project {
	return new(model.Project)
}

func GetApp(app *model.Project) *model.Project {
	return new(model.Project)
}

func GetDeployment(dpl *model.Deployment) *model.Deployment {
	return new(model.Deployment)
}
