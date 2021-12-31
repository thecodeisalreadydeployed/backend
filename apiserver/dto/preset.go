package dto

import "github.com/thecodeisalreadydeployed/model"

type CreatePresetRequest struct {
	Name     string `validate:"required"`
	Template string `validate:"required"`
}

func (req *CreatePresetRequest) ToModel() model.Preset {
	return model.Preset{
		Name:     req.Name,
		Template: req.Template,
	}
}
