package datastore

import (
	"strings"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
)

func GetDeploymentsByAppID(appID string) []model.Deployment {
	if !strings.HasPrefix(appID, "app_") {
		return []model.Deployment{}
	}

	var _data []datamodel.Deployment
	err := getDB().Table("deployments").Where(datamodel.Deployment{AppID: appID}).Scan(&_data).Error

	if err != nil {
		return []model.Deployment{}
	}

	var ret []model.Deployment
	for _, data := range _data {
		m := data.ToModel()
		ret = append(ret, m)
	}

	return ret
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
