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

	var data []model.Deployment
	err := getDB().Table("deployments").Where(datamodel.Deployment{AppID: appID}).Scan(&data).Error

	if err != nil {
		return []model.Deployment{}
	}

	return data
}

func GetDeploymentByID(deploymentID string) model.Deployment {
	if !strings.HasPrefix(deploymentID, "dpl_") {
		return model.Deployment{}
	}

	var data model.Deployment
	err := getDB().Table("deployments").Where("ID = ?", deploymentID).Scan(&data).Error

	if err != nil {
		return model.Deployment{}
	}

	return data
}
