package datastore

import (
	"github.com/thecodeisalreadydeployed/datamodel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func Init() {
	dsn := "host=localhost user=user password=password dbname=codedeploy port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&datamodel.Project{})
	database.AutoMigrate(&datamodel.App{})
	database.AutoMigrate(&datamodel.Deployment{})
	seed()
}

func getDB() *gorm.DB {
	return database
}
