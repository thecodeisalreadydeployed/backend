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

	var _data []datamodel.App
	err := getDB().Table("apps").Where(datamodel.App{ProjectID: projectID}).Scan(&_data).Error

	if err != nil {
		return []model.App{}
	}

	var ret []model.App
	for _, data := range _data {
		m := data.ToModel()
		ret = append(ret, m)
	}

	return ret
}

func GetAppByID(appID string) model.App {
	if !strings.HasPrefix(appID, "app_") {
		return model.App{}
	}

	var _data datamodel.App
	err := getDB().Table("apps").Where("ID = ?", appID).Scan(&_data).Error

	if err != nil {
		return model.App{}
	}

	return _data.ToModel()
}
