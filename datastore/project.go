package datastore

import (
	"strings"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/logger"
	"github.com/thecodeisalreadydeployed/model"
)

func GetProjectByID(projectID string) model.Project {
	if !strings.HasPrefix(projectID, "prj_") {
		return model.Project{}
	}

	var _data datamodel.Project
	err := getDB().Table("projects").Where("ID = ?", projectID).Scan(&_data).Error

	if err != nil {
		logger.Warn(err.Error())
		return model.Project{}
	}

	return _data.ToModel()
}
