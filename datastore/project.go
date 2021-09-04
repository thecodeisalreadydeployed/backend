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

func CreateProject(_p *model.Project) error {
	if !strings.HasPrefix(_p.ID, "prj_") {
		zap.L().Error(MsgProjectPrefix)
		return ErrInvalidArgument
	}
	p := datamodel.NewProjectFromModel(_p)
	err := getDB().Create(p).Error

	if err != nil {
		zap.L().Error(err.Error())
		return ErrCannotCreate
	}
	return nil
}
