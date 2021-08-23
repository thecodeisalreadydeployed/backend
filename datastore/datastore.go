package datastore

import (
	"github.com/thecodeisalreadydeployed/datamodel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := "host=localhost user=user password=password dbname=codedeploy port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = database
	DB.AutoMigrate(&datamodel.Project{})
	DB.AutoMigrate(&datamodel.App{})
	DB.AutoMigrate(&datamodel.Deployment{})
	seed()
}

func getDB() *gorm.DB {
	return DB
}
