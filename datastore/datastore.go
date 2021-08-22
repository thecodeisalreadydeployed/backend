package datastore

import (
	"fmt"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/logger"
)

func GetProjectApps(key string) []datamodel.BareApp {
	var res []datamodel.BareApp
	err := getDB().Table("apps").Where(datamodel.App{ProjectID: key}).Scan(&res).Error
	if err != nil {
		logger.Error(err.Error())
	}
	return res
}

func GetAppByID(app datamodel.App) datamodel.BareApp {
	var res datamodel.BareApp
	err := getDB().Table("apps").Where("ID = ?", app.ID).Scan(&res).Error
	if err != nil {
		logger.Error(err.Error())
	}
	return res
}

func GetAppDeployments(key string) []datamodel.BareDeployment {
	var res []datamodel.BareDeployment
	err := getDB().Table("deployments").Where(datamodel.Deployment{AppID: key}).Scan(&res).Error
	if err != nil {
		logger.Error(err.Error())
	}
	return res
}

func GetDeploymentByID(dpl datamodel.Deployment) datamodel.BareDeployment {
	var res datamodel.BareDeployment
	err := getDB().Table("deployments").Where("ID = ?", dpl.ID).Scan(&res).Error
	if err != nil {
		logger.Error(err.Error())
	}
	return res
}

func GetEvent(id string) string {
	return fmt.Sprintf("Dummy event %s.", id)
}
