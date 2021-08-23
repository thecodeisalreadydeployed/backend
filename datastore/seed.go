package datastore

import (
	"encoding/json"
	"fmt"
	faker "github.com/bxcodec/faker/v3"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/logger"
	"github.com/thecodeisalreadydeployed/model"
	"gorm.io/gorm"
	"math/rand"
)

func seed(db *gorm.DB) {
	//TODO: Updated at may come before created at

	seedProjects(db, 20)
	seedApps(db, 100)
	seedDeployments(db, 500)
}

func checkSeedExists(db *gorm.DB, name string) {
	var existing int64
	err := db.Table(name).Count(&existing).Error
	if err != nil {
		logger.Fatal(err.Error())
		return
	} else if existing > 0 {
		logger.Info(fmt.Sprintf("Table %s already seeded.", name))
		return
	}
}

func seedProjects(db *gorm.DB, size int) {
	checkSeedExists(db, "projects")

	var data []datamodel.Project
	for i := 0; i < size; i++ {
		var datum datamodel.Project
		err := faker.FakeData(&datum)
		if err != nil {
			logger.Error(err.Error())
		}
		data = append(data, datum)
	}
	db.Create(&data)
}

func seedApps(db *gorm.DB, size int) {
	checkSeedExists(db, "apps")
	var keys []string
	err := db.Table("projects").Select("ID").Scan(&keys).Error
	if err != nil {
		logger.Error(err.Error())
	}

	var data []datamodel.App
	for i := 0; i < size; i++ {
		var datum datamodel.App
		err := faker.FakeData(&datum)
		if err != nil {
			logger.Error(err.Error())
		}

		setAppForeignKey(&datum, keys)
		setAppGitSource(&datum)

		data = append(data, datum)
	}
	db.Omit("Project").Create(&data)
}

func setAppForeignKey(datum *datamodel.App, keys []string) {
	index := rand.Intn(len(keys))
	datum.ProjectID = keys[index]
}

func seedDeployments(db *gorm.DB, size int) {
	checkSeedExists(db, "deployments")

	var keys []string
	err := db.Table("apps").Select("ID").Scan(&keys).Error
	if err != nil {
		logger.Error(err.Error())
	}

	var data []datamodel.Deployment
	for i := 0; i < size; i++ {
		var datum datamodel.Deployment
		err := faker.FakeData(&datum)
		if err != nil {
			logger.Error(err.Error())
		}

		setDeploymentForeignKey(&datum, keys)
		setDeploymentGitSource(&datum)
		setDeploymentCreator(&datum)

		data = append(data, datum)
	}
	db.Omit("App").Create(&data)
}

func setDeploymentForeignKey(datum *datamodel.Deployment, keys []string) {
	index := rand.Intn(len(keys))
	datum.AppID = keys[index]
}

func setAppGitSource(datum *datamodel.App) {
	var gs model.GitSource
	err := faker.FakeData(&gs)
	if err != nil {
		logger.Error(err.Error())
	}
	res, err := json.Marshal(gs)
	datum.GitSource = string(res)
}

func setDeploymentGitSource(datum *datamodel.Deployment) {
	var gs model.GitSource
	err := faker.FakeData(&gs)
	if err != nil {
		logger.Error(err.Error())
	}
	res, err := json.Marshal(gs)
	datum.GitSource = string(res)
}

func setDeploymentCreator(datum *datamodel.Deployment) {
	var c model.Actor
	err := faker.FakeData(&c)
	if err != nil {
		logger.Error(err.Error())
	}
	res, err := json.Marshal(c)
	datum.Creator = string(res)
}
