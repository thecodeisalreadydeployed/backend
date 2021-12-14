package dto

import (
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type CreateDeploymentEventRequest struct {
	Text       string          `validate:"required"`
	ExportedAt time.Time       `validate:"required"`
	Type       model.EventType `validate:"required"`
}

func (req *CreateDeploymentEventRequest) ToModel() model.Event {
	return model.Event{
		Text:       req.Text,
		Type:       req.Type,
		ExportedAt: req.ExportedAt,
	}
}
