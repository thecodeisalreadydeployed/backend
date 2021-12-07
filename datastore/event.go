package datastore

import "gorm.io/gorm"

func GetEventByDeploymentID(DB *gorm.DB, deploymentID string) (string, error) {
	return "", nil
}
