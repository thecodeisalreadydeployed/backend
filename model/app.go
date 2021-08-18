package model

import "time"

type App struct {
	ID        string
	Name      string
	GitSource string
	CreatedAt time.Time
	UpdatedAt time.Time
}