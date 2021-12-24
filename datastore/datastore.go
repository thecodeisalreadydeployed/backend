package datastore

import (
	"fmt"
	"os"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// dsn := "host=localhost user=user password=password dbname=codedeploy port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = database
	err = DB.AutoMigrate(&datamodel.Project{})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&datamodel.App{})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&datamodel.Deployment{})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&datamodel.Event{})
	if err != nil {
		panic(err)
	}

	if util.IsDevEnvironment() {
		seed()
	}
}

func GetDB() *gorm.DB {
	return DB
}

func IsReady() bool {
	_sql, err := GetDB().DB()
	if err != nil {
		return false
	}

	err = _sql.Ping()
	return err == nil
}
