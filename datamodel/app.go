package datamodel

import "time"

type App struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	GitSource string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
