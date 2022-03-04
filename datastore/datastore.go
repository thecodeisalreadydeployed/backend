package datastore

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"log"
	"os"
	"testing"
	"time"

	"github.com/thecodeisalreadydeployed/datamodel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DataStore interface {
	CreateSeed()
	IsReady() bool

	GetAllApps() (*[]model.App, error)
	GetObservableApps() (*[]model.App, error)
	SetObservable(appID string, observable bool) error
	IsObservableApp(appID string) (bool, error)
	GetAppsByProjectID(projectID string) (*[]model.App, error)
	GetAppByID(appID string) (*model.App, error)
	SaveApp(app *model.App) (*model.App, error)
	RemoveApp(id string) error
	GetAppsByName(name string) (*[]model.App, error)

	GetPendingDeployments() (*[]model.Deployment, error)
	GetDeploymentsByAppID(appID string) (*[]model.Deployment, error)
	GetDeploymentByID(deploymentID string) (*model.Deployment, error)
	SetDeploymentState(deploymentID string, state model.DeploymentState) error
	SaveDeployment(deployment *model.Deployment) (*model.Deployment, error)
	RemoveDeployment(id string) error

	GetEventsByDeploymentID(deploymentID string) (*[]model.Event, error)
	GetEventByID(eventID string) (*model.Event, error)
	SaveEvent(event *model.Event) (*model.Event, error)
	IsValidKSUID(str string) bool

	GetAllPresets() (*[]model.Preset, error)
	GetPresetByID(presetID string) (*model.Preset, error)
	GetPresetsByName(name string) (*[]model.Preset, error)
	SavePreset(preset *model.Preset) (*model.Preset, error)
	RemovePreset(id string) error

	GetAllProjects() (*[]model.Project, error)
	GetProjectByID(id string) (*model.Project, error)
	SaveProject(project *model.Project) (*model.Project, error)
	RemoveProject(id string) error
	GetProjectsByName(name string) (*[]model.Project, error)
}

type dataStore struct {
	DB     *gorm.DB
	logger *zap.Logger
}

func NewDataStore() DataStore {
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

	var DB = database
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

	return &dataStore{
		DB:     DB,
		logger: zap.L(),
	}
}

func NewMockDataStore(gdb *gorm.DB, t *testing.T) DataStore {
	return &dataStore{
		DB:     gdb,
		logger: zaptest.NewLogger(t),
	}
}

func (d *dataStore) CreateSeed() {
	d.seedPreset()
	if util.IsDevEnvironment() {
		d.seed()
	}
}

func (d *dataStore) IsReady() bool {
	_sql, err := d.DB.DB()
	if err != nil {
		return false
	}

	err = _sql.Ping()
	return err == nil
}
