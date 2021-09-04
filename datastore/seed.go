package datastore

import (
	"encoding/json"
	"fmt"
	"math/rand"

	faker "github.com/bxcodec/faker/v3"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"go.uber.org/zap"
)

func seed() {
	//TODO: Updated at may come before created at

	seedProjects(20)
	seedApps(100)
	seedDeployments(500)
}

func seedExists(name string) bool {
	var existing int64
	err := getDB().Table(name).Count(&existing).Error
	if err != nil {
		zap.L().Error(err.Error())
		return false
	} else if existing > 0 {
		zap.L().Info(fmt.Sprintf("Table '%s' already seeded.", name))
		return true
	}
	return false
}

func seedProjects(size int) {
	if seedExists("projects") {
		return
	}

	var data []datamodel.Project
	for i := 0; i < size; i++ {
		var datum datamodel.Project
		err := faker.FakeData(&datum)
		if err != nil {
			zap.L().Error(err.Error())
		}

		datum.ID = getPrefix(datum.ID, "prj")
		data = append(data, datum)
	}
	if err := getDB().Create(&data).Error; err != nil {
		zap.L().Error("Failed to seed projects.")
	}

}

func seedApps(size int) {
	if seedExists("apps") {
		return
	}

	var keys []string
	err := getDB().Table("projects").Select("ID").Scan(&keys).Error
	if err != nil {
		zap.L().Error(err.Error())
	}

	var data []datamodel.App
	for i := 0; i < size; i++ {
		var datum datamodel.App
		err := faker.FakeData(&datum)
		if err != nil {
			zap.L().Error(err.Error())
		}

		datum.ID = getPrefix(datum.ID, "app")
		datum.ProjectID = getForeignKey(keys)
		datum.GitSource = getGitSource()

		data = append(data, datum)
	}
	if err := getDB().Omit("Project").Create(&data).Error; err != nil {
		zap.L().Error("Failed to seed apps.")
	}
}

func seedDeployments(size int) {
	if seedExists("deployments") {
		return
	}

	var keys []string
	err := getDB().Table("apps").Select("ID").Scan(&keys).Error
	if err != nil {
		zap.L().Error(err.Error())
	}

	var data []datamodel.Deployment
	for i := 0; i < size; i++ {
		var datum datamodel.Deployment
		err := faker.FakeData(&datum)
		if err != nil {
			zap.L().Error(err.Error())
		}

		datum.ID = getPrefix(datum.ID, "dpl")
		datum.AppID = getForeignKey(keys)
		datum.GitSource = getGitSource()
		datum.Creator = getCreator()
		datum.State = model.DeploymentState(getState())

		data = append(data, datum)
	}
	if err := getDB().Omit("App").Create(&data).Error; err != nil {
		zap.L().Error("Failed to seed deployments.")
	}

}

func getForeignKey(keys []string) string {
	return keys[rand.Intn(len(keys))]
}

func getGitSource() string {
	var gs model.GitSource
	err := faker.FakeData(&gs)
	if err != nil {
		zap.L().Error(err.Error())
	}
	res, err := json.Marshal(gs)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return string(res)
}

func getCreator() string {
	var c model.Actor
	err := faker.FakeData(&c)
	if err != nil {
		zap.L().Error(err.Error())
	}
	res, err := json.Marshal(c)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return string(res)
}

func getPrefix(body string, prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, body)
}

func getState() string {
	states := []string{
		"DeploymentStateQueueing",
		"DeploymentStateBuilding",
		"DeploymentStateReady",
		"DeploymentStateError",
	}
	return states[rand.Intn(4)]
}
