package datastore

import (
	"strings"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
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
