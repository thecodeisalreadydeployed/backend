package datamodel

import (
	"reflect"
	"time"

	"github.com/thecodeisalreadydeployed/model"
)

type Project struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (p *Project) ToModel() model.Project {
	return model.Project{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func NewProjectFromModel(prj *model.Project) *Project {
	return &Project{
		ID:        prj.ID,
		Name:      prj.Name,
		CreatedAt: prj.CreatedAt,
		UpdatedAt: prj.UpdatedAt,
	}
}

func ProjectStructString() []string {
	prj := Project{}
	var str []string

	e := reflect.ValueOf(&prj).Elem()

	for i := 0; i < e.NumField(); i++ {
		if IsStoredInDatabase(e.Type().Field(i)) {
			str = append(str, e.Type().Field(i).Name)
		}
	}

	return str
}
