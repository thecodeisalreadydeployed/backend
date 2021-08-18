package datamodel

import "github.com/thecodeisalreadydeployed/model"

func NewProjectFromModel(prj model.Project) Project {
	return Project{}
}

func NewAppFromModel(app model.App) App {
	return App{}
}

func NewDeploymentFromModel(dpl model.Deployment) Deployment {
	return Deployment{}
}
