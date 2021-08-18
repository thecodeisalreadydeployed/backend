package datamodel

import "github.com/thecodeisalreadydeployed/model"

func NewProjectFromModel() Project {
	return Project{}
}

func NewAppFromModel(app model.App) App {
	return App{}
}
