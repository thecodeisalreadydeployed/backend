package datastore

import (
	"strings"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
)

func GetDeploymentsByAppID(appID string) (*([]model.Deployment), error) {
	if !strings.HasPrefix(appID, "app_") {
		return nil, ErrInvalidArgument
	}

	var _data []datamodel.Deployment
	err := getDB().Table("deployments").Where(datamodel.Deployment{AppID: appID}).Scan(&_data).Error

	if err != nil {
		return nil, ErrNotFound
	}

	var _ret []model.Deployment
	for _, data := range _data {
		m := data.ToModel()
		_ret = append(_ret, m)
	}

	ret := &_ret
	return ret, nil
}

func GetDeploymentByID(deploymentID string) model.Deployment {
	if !strings.HasPrefix(deploymentID, "dpl_") {
		return model.Deployment{}
	}

	var _data datamodel.Deployment
	err := getDB().Table("deployments").Where("ID = ?", deploymentID).Scan(&_data).Error

	if err != nil {
		return model.Deployment{}
	}

	return _data.ToModel()
}
