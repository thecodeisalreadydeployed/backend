package datastore

import (
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"strings"
)

func GetProjectByID(projectID string) (*model.Project, error) {
	if !strings.HasPrefix(projectID, "prj_") {
		return nil, ErrInvalidArgument
	}

	var _data datamodel.Project

	result := getDB().First(&_data, "id = ?", projectID)

	if result.Error != nil {
		return nil, ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func CreateProject(_p *model.Project) error {
	if !strings.HasPrefix(_p.ID, "prj_") {
		return ErrInvalidArgument
	}
	p := datamodel.NewProjectFromModel(_p)
	err := getDB().Create(p).Error

	if err != nil {
		return ErrCannotCreate
	}
	return nil
}
