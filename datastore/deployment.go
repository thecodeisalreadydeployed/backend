package datastore

import (
	"gorm.io/gorm"
	"strings"

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

func SetDeploymentState(deploymentID string, state model.DeploymentState) error {
	dpl, getDeploymentErr := GetDeploymentByID(DB, deploymentID)
	if getDeploymentErr != nil {
		return getDeploymentErr
	}

	err := DB.Model(&dpl).Update("state", state).Error
	if err != nil {
		return err
	}

	return nil
}
