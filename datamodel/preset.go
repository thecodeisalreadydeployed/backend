package datamodel

import "github.com/thecodeisalreadydeployed/model"

type Preset struct {
	ID   string `gorm:"primaryKey"`
	Name string

	// Stored in base64 encoding.
	Template string
}

func (p *Preset) ToModel() model.Preset {
	return model.Preset{
		ID:       p.ID,
		Name:     p.Name,
		Template: model.GetDecodedString(p.Template),
	}
}

func NewPresetFromModel(p *model.Preset) *Preset {
	return &Preset{
		ID:       p.ID,
		Name:     p.Name,
		Template: model.GetEncodedString(p.Template),
	}
}
