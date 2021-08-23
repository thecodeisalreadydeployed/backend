package datastore

import (
	"strings"

	"github.com/thecodeisalreadydeployed/logger"
	"github.com/thecodeisalreadydeployed/model"
)

func GetProjectByID(projectID string) model.Project {
	if !strings.HasPrefix(projectID, "prj_") {
		return model.Project{}
	}

	var data model.Project
	err := getDB().Table("projects").Where("ID = ?", projectID).Scan(&data).Error

	if err != nil {
		logger.Warn(err.Error())
		return model.Project{}
	}

	return data
}
