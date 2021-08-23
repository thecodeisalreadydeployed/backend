package datastore

import (
	"strings"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/logger"
	"github.com/thecodeisalreadydeployed/model"
)

func GetProjectByID(projectID string) (*model.Project, error) {
	if !strings.HasPrefix(projectID, "prj_") {
		return nil, ErrInvalidArgument
	}

	var _data datamodel.Project
	err := getDB().Table("projects").Where("ID = ?", projectID).Scan(&_data).Error

	if err != nil {
		logger.Warn(err.Error())
		return nil, ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}
