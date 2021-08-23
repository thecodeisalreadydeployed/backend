package datastore

import (
	"strings"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
)

func GetAppsByProjectID(projectID string) []model.App {
	if !strings.HasPrefix(projectID, "prj_") {
		return []model.App{}
	}

	var data []model.App
	err := getDB().Table("apps").Where(datamodel.App{ProjectID: projectID}).Scan(&data).Error

	if err != nil {
		return []model.App{}
	}

	return data
}

func GetAppByID(appID string) model.App {
	if !strings.HasPrefix(appID, "app_") {
		return model.App{}
	}

	var data model.App
	err := getDB().Table("apps").Where("ID = ?", appID).Scan(&data).Error

	if err != nil {
		return model.App{}
	}

	return data
}
