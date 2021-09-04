package datastore

import (
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
	"strings"
)

func GetProjectByID(projectID string) (*model.Project, error) {
	if !strings.HasPrefix(projectID, "prj_") {
		zap.L().Error(MsgProjectPrefix)
		return nil, ErrInvalidArgument
	}

	var _data datamodel.Project

	err := getDB().First(&_data, "id = ?", projectID).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func SaveProject(_p *model.Project) error {
	if !strings.HasPrefix(_p.ID, "prj_") {
		zap.L().Error(MsgProjectPrefix)
		return ErrInvalidArgument
	}
	p := datamodel.NewProjectFromModel(_p)
	err := getDB().Save(p).Error

	if err != nil {
		zap.L().Error(err.Error())
		return ErrCannotSave
	}
	return nil
}

func RemoveProject(id string) error {
	if !strings.HasPrefix(id, "prj_") {
		zap.L().Error(MsgProjectPrefix)
		return ErrInvalidArgument
	}
	var p datamodel.Project
	err := getDB().Table("projects").Where(datamodel.Project{ID: id}).First(&p).Error
	if err != nil {
		zap.L().Error(err.Error())
		return ErrNotFound
	}
	if err := getDB().Delete(&p).Error; err != nil {
		zap.L().Error(err.Error())
		return ErrCannotDelete
	}
	return nil
}
