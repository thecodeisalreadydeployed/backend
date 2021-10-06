package dto

import (
	"github.com/thecodeisalreadydeployed/model"
	"time"
)

type CreateProjectRequest struct {
	Name string `validate:"required"`
}

func (req *CreateProjectRequest) ToModel() model.Project {
	return model.Project{
		ID:        model.GenerateProjectID(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
