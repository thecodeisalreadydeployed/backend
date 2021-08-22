package datastore

import (
	"fmt"
	faker "github.com/bxcodec/faker/v3"
	"github.com/thecodeisalreadydeployed/datamodel"
	"gorm.io/gorm"
	"math/rand"
)

func seed(db *gorm.DB) {
	//TODO: Updated at may come before created at

	seedProjects(db, 20)
	seedApps(db, 100)
	seedDeployments(db, 500)
}

func checkSeedExists(db *gorm.DB, tablename string) {
	var existing int64
	err := db.Table(tablename).Count(&existing).Error
	if err != nil {
		fmt.Println(err)
		return
	} else if existing > 0 {
		fmt.Printf("Table %s already seeded.", tablename)
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
			fmt.Println(err)
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
		fmt.Println(err)
	}

	var data []datamodel.App
	for i := 0; i < size; i++ {
		var datum datamodel.App
		err := faker.FakeData(&datum)
		if err != nil {
			fmt.Println(err)
		}

		setAppForeignKey(&datum, keys)

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
		fmt.Println(err)
	}

	var data []datamodel.Deployment
	for i := 0; i < size; i++ {
		var datum datamodel.Deployment
		err := faker.FakeData(&datum)
		if err != nil {
			fmt.Println(err)
		}

		setDeploymentForeignKey(&datum, keys)

		data = append(data, datum)
	}
	db.Omit("App").Create(&data)
}

func setDeploymentForeignKey(datum *datamodel.Deployment, keys []string) {
	index := rand.Intn(len(keys))
	datum.AppID = keys[index]
}

