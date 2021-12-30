package datastore

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"go.uber.org/zap"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/errutil"
	"github.com/thecodeisalreadydeployed/model"
)

func GetDeploymentsByAppID(DB *gorm.DB, appID string) (*[]model.Deployment, error) {
	if !strings.HasPrefix(appID, "app_") {
		zap.L().Error(MsgAppPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data []datamodel.Deployment
	err := DB.Table("deployments").Where(datamodel.Deployment{AppID: appID}).Scan(&_data).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	var _ret []model.Deployment
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
}

func GetDeploymentByID(DB *gorm.DB, deploymentID string) (*model.Deployment, error) {
	if !strings.HasPrefix(deploymentID, "dpl_") {
		zap.L().Error(MsgDeploymentPrefix)
		return nil, errutil.ErrInvalidArgument
	}

	var _data datamodel.Deployment
	err := DB.First(&_data, "id = ?", deploymentID).Error

	if err != nil {
		zap.L().Error(err.Error())
		return nil, errutil.ErrNotFound
	}

	ret := _data.ToModel()
	return &ret, nil
}

func SetDeploymentState(DB *gorm.DB, deploymentID string, state model.DeploymentState) error {
	err := DB.Table("deployments").
		Where(datamodel.Deployment{ID: deploymentID}).
		Update("state", state).
		Error
	if err != nil {
		return err
	}

	return nil
}

func SaveDeployment(DB *gorm.DB, deployment *model.Deployment) (*model.Deployment, error) {
	if deployment.ID != "" {
		if !strings.HasPrefix(deployment.ID, "dpl_") {
			return nil, errutil.ErrInvalidArgument
		}
	} else {
		deployment.ID = model.GenerateDeploymentID()
	}

	a := datamodel.NewDeploymentFromModel(deployment)
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
	return GetDeploymentByID(DB, deployment.ID)
}
