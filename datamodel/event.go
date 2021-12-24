package datamodel

import (
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type Event struct {
	ID           string `gorm:"primaryKey"`
	DeploymentID string
	Text         string
	Type         model.EventType
	ExportedAt   time.Time
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (e *Event) ToModel() model.Event {
	return model.Event{
		ID:           e.ID,
		DeploymentID: e.DeploymentID,
		Text:         e.Text,
		Type:         e.Type,
		ExportedAt:   e.ExportedAt,
		CreatedAt:    e.CreatedAt,
	}
}

func NewEventFromModel(event *model.Event) *Event {
	return &Event{
		ID:           event.ID,
		DeploymentID: event.DeploymentID,
		Text:         event.Text,
		Type:         event.Type,
		ExportedAt:   event.ExportedAt,
		CreatedAt:    event.CreatedAt,
	}
}
