package datastore

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func ExpectVersionQuery(mock sqlmock.Sqlmock) {
	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mock.NewRows([]string{"version"}).FromCSVString("1"),
	)
}

func OpenGormDB(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
}

func getDeploymentRows() *sqlmock.Rows {
	return sqlmock.NewRows(datamodel.DeploymentStructString()).
		AddRow(
			"dpl_test",
			"app_test",
			model.GetCreatorString(model.Actor{}),
			"dummy_meta",
			model.GetGitSourceString(model.GitSource{}),
			time.Unix(0, 0),
			time.Unix(0, 0),
			time.Unix(0, 0),
			model.GetBuildConfigurationString(model.BuildConfiguration{}),
			time.Unix(0, 0),
			time.Unix(0, 0),
			model.DeploymentStateReady,
		)
}

func getAppRows() *sqlmock.Rows {
	return sqlmock.NewRows(datamodel.AppStructString()).
		AddRow(
			"app_test",
			"prj_test",
			"Best App",
			model.GetGitSourceString(model.GitSource{}),
			time.Unix(0, 0),
			time.Unix(0, 0),
			model.GetBuildConfigurationString(model.BuildConfiguration{}),
			true,
		)
}

func getProjectRows() *sqlmock.Rows {
	return sqlmock.NewRows(datamodel.ProjectStructString()).
		AddRow("prj_test", "Best Project", time.Unix(0, 0), time.Unix(0, 0))
}

func getExpectedDeployment() *model.Deployment {
	return &model.Deployment{
		ID:                 "dpl_test",
		AppID:              "app_test",
		Creator:            model.Actor{},
		Meta:               "dummy_meta",
		GitSource:          model.GitSource{},
		BuiltAt:            time.Unix(0, 0),
		CommittedAt:        time.Unix(0, 0),
		DeployedAt:         time.Unix(0, 0),
		BuildConfiguration: model.BuildConfiguration{},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		State:              model.DeploymentStateReady,
	}
}

func getExpectedApp() *model.App {
	return &model.App{
		ID:                 "app_test",
		ProjectID:          "prj_test",
		Name:               "Best App",
		GitSource:          model.GitSource{},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		BuildConfiguration: model.BuildConfiguration{},
		Observable:         true,
	}
}

func getExpectedProject() *model.Project {
	return &model.Project{
		ID:        "prj_test",
		Name:      "Best Project",
		CreatedAt: time.Unix(0, 0),
		UpdatedAt: time.Unix(0, 0),
	}
}
