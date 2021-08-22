package datastore

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
)

func GetDB() *gorm.DB {
	dsn := "host=localhost user=user password=password dbname=codedeploy port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func InitDB(db *gorm.DB) {
	createTable(db, &datamodel.Project{})
	createTable(db, &datamodel.App{})
	createTable(db, &datamodel.Deployment{})

	seed(db)
}

func createTable(db *gorm.DB, i interface{}) {
	if !db.Migrator().HasTable(i) {
		err := db.Migrator().CreateTable(i)
		if err != nil {
			panic(err)
		}
	} else {
		name := reflect.TypeOf(i).Elem().Name()
		fmt.Printf("Table %s already created.\n", name)
	}
}

//func GetAllProjects(p *model.Project) *[]model.Project {
//
//}
//
//func GetAllAppsFromProject(app *model.App) *[]model.App {
//
//}
//
//func GetAllDeploymentsFromApp(dpl *model.Deployment) *[]model.Deployment {
//
//}
func GetProject(db *gorm.DB, p *datamodel.Project) *datamodel.Project {
	GetDB()
	var result datamodel.Project
	db.First(&result)
	return &result
}

func GetApp(app *model.App) *model.App {
	return new(model.App)
}

func GetDeployment(dpl *model.Deployment) *model.Deployment {
	return new(model.Deployment)
}

func GetEvent(id string) string {
	return "Dummy event."
}
