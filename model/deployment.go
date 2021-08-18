package model

import "time"

type Deployment struct {
	ID        string
	Name      string
	Creator   string
	Meta      string
	GitBranch string
	BuildedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	State     DeploymentState
}