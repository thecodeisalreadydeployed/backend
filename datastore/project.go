package datastore

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
)

func (d *dataStore) GetAllProjects() (*[]model.Project, error) {
	var _data []datamodel.Project
	err := d.DB.Table("projects").Scan(&_data).Error
	if err != nil {
		d.logger.Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	var _ret []model.Project
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
}

func (d *dataStore) GetProjectByID(id string) (*model.Project, error) {
	if !strings.HasPrefix(id, "prj-") {
		d.logger.Error(MsgProjectPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data datamodel.Project

	err := d.DB.Table("projects").First(&_data, "id = ?", id).Error

	if err != nil {
		d.logger.Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func (d *dataStore) SaveProject(project *model.Project) (*model.Project, error) {
	if project.ID == "" {
		project.ID = model.GenerateProjectID()
	}
	if !strings.HasPrefix(project.ID, "prj-") {
		d.logger.Error(MsgProjectPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	p := datamodel.NewProjectFromModel(project)
	err := d.DB.Save(p).Error

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

	return d.GetProjectByID(project.ID)
}

func (d *dataStore) RemoveProject(id string) error {
	if !strings.HasPrefix(id, "prj-") {
		d.logger.Error(MsgProjectPrefix)
		return errutil.ErrInvalidArgument
	}
	var p datamodel.Project

	err := d.DB.Table("projects").Where(datamodel.Project{ID: id}).First(&p).Error
	if err != nil {
		d.logger.Error(err.Error())
		return errutil.ErrNotFound
	}
	if err := d.DB.Delete(&p).Error; err != nil {
		d.logger.Error(err.Error())
		return errutil.ErrUnknown
	}
	return nil
}

func (d *dataStore) GetProjectsByName(name string) (*[]model.Project, error) {
	var _data []datamodel.Project

	err := d.DB.Table("projects").Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).Scan(&_data).Error

	if err != nil {
		d.logger.Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	var _ret []model.Project
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
}
