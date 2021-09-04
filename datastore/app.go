package datastore

import (
	"go.uber.org/zap"
	"strings"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
)

func GetAllApps() (*[]model.App, error) {
	var _data []datamodel.App
	err := getDB().Table("apps").Scan(&_data).Error
	if err != nil {
		zap.L().Error(err.Error())
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

func GetAppsByProjectID(projectID string) (*[]model.App, error) {
	if !strings.HasPrefix(projectID, "prj_") {
		zap.L().Error(MsgProjectPrefix)
		return nil, ErrInvalidArgument
	}

	var _data []datamodel.App
	err := getDB().Table("apps").Where(datamodel.App{ProjectID: projectID}).Scan(&_data).Error

	if err != nil {
		zap.L().Error(err.Error())
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
		zap.L().Error(MsgAppPrefix)
		return nil, ErrInvalidArgument
	}

	var _data datamodel.App
	err := getDB().First(&_data, "id = ?", appID).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func SaveApp(_a *model.App) error {
	if !strings.HasPrefix(_a.ID, "app_") {
		zap.L().Error(MsgAppPrefix)
		return ErrInvalidArgument
	}
	a := datamodel.NewAppFromModel(_a)
	err := getDB().Save(a).Error

	if err != nil {
		zap.L().Error(err.Error())
		return ErrCannotSave
	}
	return nil
}

func RemoveApp(id string) error {
	if !strings.HasPrefix(id, "app_") {
		zap.L().Error(MsgAppPrefix)
		return ErrInvalidArgument
	}
	var a datamodel.App
	err := getDB().Table("apps").Where(datamodel.App{ID: id}).First(&a).Error
	if err != nil {
		zap.L().Error(err.Error())
		return ErrNotFound
	}
	if err := getDB().Delete(&a).Error; err != nil {
		zap.L().Error(err.Error())
		return ErrCannotDelete
	}
	return nil
}
