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
	Finish       bool      `json:"-"`
}

func GenerateEventID(exportedAt time.Time) string {
	id, err := ksuid.NewRandomWithTime(exportedAt)
	if err != nil {
		panic("cannot generate KSUID")
	}

	return id.String()
}
