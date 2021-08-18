package datamodel

import "time"

type Project struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewProjectFromModel() Project {
	return Project{}
}

func (p *Project) toModel() {}
