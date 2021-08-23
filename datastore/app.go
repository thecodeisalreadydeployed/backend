package datastore

import (
	"strings"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
)

func GetAppsByProjectID(projectID string) (*[]model.App, error) {
	if !strings.HasPrefix(projectID, "prj_") {
		return nil, ErrInvalidArgument
	}

	var _data []datamodel.App
	err := getDB().Table("apps").Where(datamodel.App{ProjectID: projectID}).Scan(&_data).Error

	if err != nil {
		return nil, ErrNotFound
	}

	var _ret []model.App
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
}

func GetAppByID(appID string) (*model.App, error) {
	if !strings.HasPrefix(appID, "app_") {
		return nil, ErrInvalidArgument
	}

	var _data datamodel.App
	err := getDB().First(&_data, "id = ?", appID)

	if err != nil {
		return nil, ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}
