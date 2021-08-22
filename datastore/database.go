package datastore

import (
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	//TODO: Move to environment variable
	dsn := "host=localhost user=user password=password dbname=codedeploy port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error(err.Error())
	}
	return db
}

func InitDB(db *gorm.DB) {
	createTable(db, &datamodel.Project{})
	createTable(db, &datamodel.App{})
	createTable(db, &datamodel.Deployment{})

	seed(db)
}
