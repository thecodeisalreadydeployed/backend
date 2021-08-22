package datastore

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/logger"
)

func GetProject(p datamodel.Project) []datamodel.Project {
	var res []datamodel.Project
	err := getDB().Table("projects").Where(p).Scan(&res).Error
	if err != nil {
		logger.Warn(err.Error())
	}
	return res
}

func GetProjectApps(key string) []datamodel.App {
	var res []datamodel.App
	err := getDB().Table("apps").Where(datamodel.App{ProjectID: key}).Scan(&res).Error
	if err != nil {
		logger.Error(err.Error())
	}
	return res
}

func GetApp(app datamodel.App) []datamodel.App {
	var res []datamodel.App
	err := getDB().Table("apps").Where(app).Scan(&res).Error
	if err != nil {
		logger.Error(err.Error())
	}
	return res
}

func GetAppDeployments(key string) []datamodel.Deployment {
	var res []datamodel.Deployment
	err := getDB().Table("deployments").Where(datamodel.Deployment{AppID: key}).Scan(&res).Error
	if err != nil {
		logger.Error(err.Error())
	}
	return res
}

func GetDeployment(dpl datamodel.Deployment) []datamodel.Deployment {
	var res []datamodel.Deployment
	err := getDB().Table("apps").Where(dpl).Scan(&res).Error
	if err != nil {
		logger.Error(err.Error())
	}
	return res
}

func GetEvent(id string) string {
	return fmt.Sprintf("Dummy event %s.", id)
}
