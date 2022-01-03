package datastore

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/model"
	"log"
	"os"
	"sync"
	"time"

	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var appChan = make(chan *model.App)
var observables sync.Map

func GetAppChannel() chan *model.App {
	return appChan
}

func GetObservables() *sync.Map {
	return &observables
}

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

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
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

	err = DB.AutoMigrate(&datamodel.Preset{})
	if err != nil {
		panic(err)
	}

	seedPreset()
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
