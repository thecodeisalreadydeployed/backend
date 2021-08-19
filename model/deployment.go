package model

import "time"

type Deployment struct {
	ID        string
	Name      string
	Creator   Actor
	Meta      string
	GitSource GitSource
	BuildedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	State     DeploymentState
}
