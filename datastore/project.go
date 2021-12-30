package datastore

import (
	"strings"

	"gorm.io/gorm"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

func GetAllProjects(DB *gorm.DB) (*[]model.Project, error) {
	var _data []datamodel.Project
	err := DB.Table("projects").Scan(&_data).Error
	if err != nil {
		zap.L().Error(err.Error())
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

func GetProjectByID(DB *gorm.DB, id string) (*model.Project, error) {
	if !strings.HasPrefix(id, "prj-") {
		zap.L().Error(MsgProjectPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data datamodel.Project

	err := DB.Table("projects").First(&_data, "id = ?", id).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func SaveProject(DB *gorm.DB, project *model.Project) (*model.Project, error) {
	if project.ID != "" {
		if !strings.HasPrefix(project.ID, "prj-") {
			zap.L().Error(MsgProjectPrefix)
			return nil, errutil.ErrInvalidArgument
		}
	} else {
		project.ID = model.GenerateProjectID()
	}
	p := datamodel.NewProjectFromModel(project)
	err := DB.Save(p).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrUnknown
	}

	return GetProjectByID(DB, project.ID)
}

func RemoveProject(DB *gorm.DB, id string) error {
	if !strings.HasPrefix(id, "prj-") {
		zap.L().Error(MsgProjectPrefix)
		return errutil.ErrInvalidArgument
	}
	var p datamodel.Project
	err := DB.Table("projects").Where(datamodel.Project{ID: id}).First(&p).Error
	if err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrNotFound
	}
	if err := DB.Delete(&p).Error; err != nil {
		zap.L().Error(err.Error())
		return errutil.ErrUnknown
	}
	return nil
}

func GetProjectsByName(DB *gorm.DB, name string) (*[]model.Project, error) {
	var _data []datamodel.Project

	err := DB.Table("projects").Where(datamodel.Project{Name: name}).Scan(&_data).Error

	if err != nil {
		zap.L().Error(err.Error())
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
