package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

type EventType string

const (
	INFO  EventType = "INFO"
	DEBUG EventType = "DEBUG"
	ERROR EventType = "ERROR"
)

type Event struct {
	ID           string    `json:"id"`
	DeploymentID string    `json:"deploymentID"`
	Text         string    `json:"text"`
	Type         EventType `json:"type"`
	CreatedAt    time.Time `json:"createdAt"`
	ExportedAt   time.Time `json:"exportedAt"`
}

func GenerateEventID() string {
	id := ksuid.New()
	return id.String()
}
