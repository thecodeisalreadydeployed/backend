package dto

import (
	"encoding/json"
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type CreateDeploymentEventRequest struct {
	Text       string          `validate:"required"`
	ExportedAt time.Time       `validate:"required"`
	Type       model.EventType `validate:"required"`
	Finish     bool            `validate:"required"`
}

type KanikoLog struct {
	Level string
	Msg   string
}

func (req *CreateDeploymentEventRequest) ToModel() model.Event {
	var kanikoLog KanikoLog
	err := json.Unmarshal([]byte(req.Text), &kanikoLog)
	if err != nil {
		return model.Event{
			Text:       req.Text,
			Type:       req.Type,
			ExportedAt: req.ExportedAt,
		}
	}

	var logLevel model.EventType
	switch kanikoLog.Level {
	case "debug":
		logLevel = model.DEBUG
	case "info":
		logLevel = model.INFO
	case "warning":
		logLevel = model.INFO
	default:
		logLevel = model.ERROR
	}

	return model.Event{
		Text:       kanikoLog.Msg,
		Type:       logLevel,
		ExportedAt: req.ExportedAt,
	}
}
