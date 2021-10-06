package datastore

import (
	"errors"
	"strings"

	"go.uber.org/zap"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	"gorm.io/gorm"
)

func GetAllApps() (*[]model.App, error) {
	var _data []datamodel.App
	err := getDB().Table("apps").Scan(&_data).Error
	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	var _ret []model.App
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
}

func GetObservableApps() (*[]model.App, error) {
	var _data []datamodel.App
	err := getDB().Table("apps").Where("observable = ?", true).Scan(&_data).Error
	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
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
		return nil, errutil.ErrInvalidArgument
	}

	var _data []datamodel.App
	err := getDB().Table("apps").Where(datamodel.App{ProjectID: projectID}).Scan(&_data).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
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
		return nil, errutil.ErrInvalidArgument
	}

	var _data datamodel.App
	err := getDB().First(&_data, "id = ?", appID).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func SaveApp(_a *model.App) error {
	if !strings.HasPrefix(_a.ID, "app_") {
		zap.L().Error(MsgAppPrefix)
		return errutil.ErrInvalidArgument
	}
	a := datamodel.NewAppFromModel(_a)
	err := getDB().Save(a).Error

	if err != nil {
		zap.L().Error(err.Error())

		if errors.Is(err, gorm.ErrInvalidField) || errors.Is(err, gorm.ErrInvalidData) {
			return errutil.ErrInvalidArgument
		}

		if errors.Is(err, gorm.ErrMissingWhereClause) {
			return errutil.ErrFailedPrecondition
		}

		return errutil.ErrUnknown
	}
	return nil
}

func RemoveApp(id string) error {
	if !strings.HasPrefix(id, "app_") {
		zap.L().Error(MsgAppPrefix)
		return errutil.ErrInvalidArgument
	}
	var a datamodel.App
	err := getDB().Table("apps").Where(datamodel.App{ID: id}).First(&a).Error
	if err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrNotFound
	}
	if err := getDB().Delete(&a).Error; err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrUnknown
	}
	return nil
}

func GetAppsByName(name string) (*[]model.App, error) {
	var _data []datamodel.App

	err := getDB().Table("apps").Where(datamodel.App{Name: name}).Scan(&_data).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	var _ret []model.App
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
}
