package datastore

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	"gorm.io/gorm"
)

func (d *dataStore) GetAllApps() (*[]model.App, error) {
	var _data []datamodel.App
	err := d.DB.Table("apps").Scan(&_data).Error
	if err != nil {
		d.logger.Error(err.Error())
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

func (d *dataStore) GetObservableApps() (*[]model.App, error) {
	var _data []datamodel.App
	err := d.DB.Table("apps").Where("observable = ?", true).Scan(&_data).Error
	if err != nil {
		d.logger.Error(err.Error())
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

func (d *dataStore) SetObservable(appID string, observable bool) error {
	var app datamodel.App
	err := d.DB.First(&app, "id = ?", appID).Error
	if err != nil {
		d.logger.Error(err.Error())
		return errutil.ErrNotFound
	}

	app.Observable = observable
	err = d.DB.Save(&app).Error
	if err != nil {
		d.logger.Error(err.Error())
		return errutil.ErrNotFound
	}

	return nil
}

func (d *dataStore) IsObservableApp(appID string) (bool, error) {
	var observable bool
	err := d.DB.Table("apps").
		Where(datamodel.App{ID: appID}).
		Select("Observable").
		Scan(&observable).Error

	if err != nil {
		d.logger.Error(err.Error())
		return false, errutil.ErrNotFound
	}
	return observable, nil
}

func (d *dataStore) GetAppsByProjectID(projectID string) (*[]model.App, error) {
	if !strings.HasPrefix(projectID, "prj-") {
		d.logger.Error(MsgProjectPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data []datamodel.App
	err := d.DB.Table("apps").Where(datamodel.App{ProjectID: projectID}).Scan(&_data).Error

	if err != nil {
		d.logger.Error(err.Error())
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

func (d *dataStore) GetAppByID(appID string) (*model.App, error) {
	if !strings.HasPrefix(appID, "app-") {
		d.logger.Error(MsgAppPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data datamodel.App
	err := d.DB.First(&_data, "id = ?", appID).Error

	if err != nil {
		d.logger.Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func (d *dataStore) SaveApp(app *model.App) (*model.App, error) {
	if app.ID == "" {
		app.ID = model.GenerateAppID()
	}
	if !strings.HasPrefix(app.ID, "app-") {
		d.logger.Error(MsgAppPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	a := datamodel.NewAppFromModel(app)
	err := d.DB.Save(a).Error

	if err != nil {
		d.logger.Error(err.Error())

		if errors.Is(err, gorm.ErrInvalidField) || errors.Is(err, gorm.ErrInvalidData) {
			return nil, errutil.ErrInvalidArgument
		}

		if errors.Is(err, gorm.ErrMissingWhereClause) {
			return nil, errutil.ErrFailedPrecondition
		}

		return nil, errutil.ErrUnknown
	}

	return d.GetAppByID(app.ID)
}

func (d *dataStore) RemoveApp(id string) error {
	if !strings.HasPrefix(id, "app-") {
		d.logger.Error(MsgAppPrefix)
		return errutil.ErrInvalidArgument
	}
	var a datamodel.App
	err := d.DB.Table("apps").Where(datamodel.App{ID: id}).First(&a).Error
	if err != nil {
		d.logger.Error(err.Error())
		return errutil.ErrNotFound
	}
	if err := d.DB.Delete(&a).Error; err != nil {
		d.logger.Error(err.Error())
		return errutil.ErrUnknown
	}
	return nil
}

func (d *dataStore) GetAppsByName(name string) (*[]model.App, error) {
	var _data []datamodel.App

	err := d.DB.Table("apps").Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).Scan(&_data).Error

	if err != nil {
		d.logger.Error(err.Error())
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
