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

func GetAllApps(DB *gorm.DB) (*[]model.App, error) {
	var _data []datamodel.App
	err := DB.Table("apps").Scan(&_data).Error
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

func GetObservableApps(DB *gorm.DB) (*[]model.App, error) {
	var _data []datamodel.App
	err := DB.Table("apps").Where("observable = ?", true).Scan(&_data).Error
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

func GetAppsByProjectID(DB *gorm.DB, projectID string) (*[]model.App, error) {
	if !strings.HasPrefix(projectID, "prj_") {
		zap.L().Error(MsgProjectPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data []datamodel.App
	err := DB.Table("apps").Where(datamodel.App{ProjectID: projectID}).Scan(&_data).Error

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

func GetAppByID(DB *gorm.DB, appID string) (*model.App, error) {
	if !strings.HasPrefix(appID, "app_") {
		zap.L().Error(MsgAppPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data datamodel.App
	err := DB.First(&_data, "id = ?", appID).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func SaveApp(DB *gorm.DB, app *model.App) (*model.App, error) {
	if app.ID != "" {
		if !strings.HasPrefix(app.ID, "app_") {
			zap.L().Error(MsgAppPrefix)
			return nil, errutil.ErrInvalidArgument
		}
	} else {
		app.ID = model.GenerateAppID()
	}

	a := datamodel.NewAppFromModel(app)
	err := DB.Save(a).Error

	if err != nil {
		zap.L().Error(err.Error())

		if errors.Is(err, gorm.ErrInvalidField) || errors.Is(err, gorm.ErrInvalidData) {
			return nil, errutil.ErrInvalidArgument
		}

		if errors.Is(err, gorm.ErrMissingWhereClause) {
			return nil, errutil.ErrFailedPrecondition
		}

		return nil, errutil.ErrUnknown
	}
	return GetAppByID(DB, app.ID)
}

func RemoveApp(DB *gorm.DB, id string) error {
	if !strings.HasPrefix(id, "app_") {
		zap.L().Error(MsgAppPrefix)
		return errutil.ErrInvalidArgument
	}
	var a datamodel.App
	err := DB.Table("apps").Where(datamodel.App{ID: id}).First(&a).Error
	if err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrNotFound
	}
	if err := DB.Delete(&a).Error; err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrUnknown
	}
	return nil
}

func GetAppsByName(DB *gorm.DB, name string) (*[]model.App, error) {
	var _data []datamodel.App

	err := DB.Table("apps").Where(datamodel.App{Name: name}).Scan(&_data).Error

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
