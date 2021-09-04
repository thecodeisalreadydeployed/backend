package dto

import "github.com/thecodeisalreadydeployed/model"

type CreateProjectRequest struct {
	Name string `validate:"required"`
}

func (req *CreateProjectRequest) ToModel() model.Project {
	return model.Project{
		Name: req.Name,
	}
}
