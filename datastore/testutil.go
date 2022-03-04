package datastore

import (
	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thecodeisalreadydeployed/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func DeploymentStructString() []string {
	return []string{
		"ID",
		"AppID",
		"Creator",
		"Meta",
		"GitSource",
		"BuiltAt",
		"CommittedAt",
		"DeployedAt",
		"BuildConfiguration",
		"CreatedAt",
		"UpdatedAt",
		"State",
	}
}

func AppStructString() []string {
	return []string{
		"ID",
		"ProjectID",
		"Name",
		"GitSource",
		"CreatedAt",
		"UpdatedAt",
		"BuildConfiguration",
		"Observable",
		"FetchInterval",
	}
}

func ProjectStructString() []string {
	return []string{"ID", "Name", "CreatedAt", "UpdatedAt"}
}

func EventStructString() []string {
	return []string{
		"ID",
		"DeploymentID",
		"Text",
		"Type",
		"ExportedAt",
		"CreatedAt",
	}
}

func PresetStructString() []string {
	return []string{"ID", "Name", "Template"}
}

func GetDeploymentRows() *sqlmock.Rows {
	return sqlmock.NewRows(DeploymentStructString()).
		AddRow(
			"dpl-test",
			"app-test",
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

func GetAppRows() *sqlmock.Rows {
	return sqlmock.NewRows(AppStructString()).
		AddRow(
			"app-test",
			"prj-test",
			"Best App",
			model.GetGitSourceString(model.GitSource{
				CommitSHA:        "a",
				CommitMessage:    "a",
				CommitAuthorName: "a",
				RepositoryURL:    "a",
				Branch:           "a",
			}),
			time.Unix(0, 0),
			time.Unix(0, 0),
			model.GetBuildConfigurationString(model.BuildConfiguration{}),
			true,
			0,
		)
}

func GetProjectRows() *sqlmock.Rows {
	return sqlmock.NewRows(ProjectStructString()).
		AddRow("prj-test", "Best Project", time.Unix(0, 0), time.Unix(0, 0))
}

func GetEventRows() *sqlmock.Rows {
	return sqlmock.NewRows(EventStructString()).
		AddRow(
			"abcdefghijklmnopqrstuvwxyz0",
			"dpl-test",
			"Downloading dependencies (1/20)",
			model.INFO,
			time.Unix(0, 0),
			time.Unix(0, 0))
}

func GetPresetRows() *sqlmock.Rows {
	return sqlmock.NewRows(PresetStructString()).AddRow("pst-test", "My Preset", "UlVOIGVjaG8gaGVsbG8=")
}

func GetExpectedDeployment() *model.Deployment {
	return &model.Deployment{
		ID:                 "dpl-test",
		AppID:              "app-test",
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

func GetExpectedApp() *model.App {
	return &model.App{
		ID:        "app-test",
		ProjectID: "prj-test",
		Name:      "Best App",
		GitSource: model.GitSource{
			CommitSHA:        "a",
			CommitMessage:    "a",
			CommitAuthorName: "a",
			RepositoryURL:    "a",
			Branch:           "a",
		},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		BuildConfiguration: model.BuildConfiguration{},
		Observable:         true,
		FetchInterval:      0,
	}
}

func GetExpectedProject() *model.Project {
	return &model.Project{
		ID:        "prj-test",
		Name:      "Best Project",
		CreatedAt: time.Unix(0, 0),
		UpdatedAt: time.Unix(0, 0),
	}
}

func GetExpectedEvent() *model.Event {
	return &model.Event{
		ID:           "abcdefghijklmnopqrstuvwxyz0",
		DeploymentID: "dpl-test",
		Text:         "Downloading dependencies (1/20)",
		Type:         model.INFO,
		CreatedAt:    time.Unix(0, 0),
		ExportedAt:   time.Unix(0, 0),
	}
}

func GetExpectedPreset() *model.Preset {
	return &model.Preset{
		ID:       "pst-test",
		Name:     "My Preset",
		Template: "RUN echo hello",
	}
}
