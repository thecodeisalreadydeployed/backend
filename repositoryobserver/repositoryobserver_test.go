package repositoryobserver

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/model"
	mock_workloadcontroller "github.com/thecodeisalreadydeployed/workloadcontroller/v2/mock"
	"go.uber.org/zap/zaptest"

	"github.com/stretchr/testify/assert"
)

func TestCheckChanges(t *testing.T) {
	logger := zaptest.NewLogger(t)
	observer := NewRepositoryObserver(nil, nil, nil)
	changeString, duration := observer.CheckChanges(
		logger,
		"https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		"main",
		"37e8e4d20d889924780f2373453a246591b6b11a",
	)

	assert.Equal(t, "5da29979c5ef986dc8ec6aa603e0862310abc96e", *changeString)
	assert.Equal(t, 19*time.Minute+57*time.Second, duration)

	changeString, duration = observer.CheckChanges(
		logger,
		"https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		"main",
		"5da29979c5ef986dc8ec6aa603e0862310abc96e",
	)

	assert.Nil(t, changeString)
	assert.Equal(t, 19*time.Minute+57*time.Second, duration)

	changeString, duration = observer.CheckChanges(
		logger,
		"https://github.com/thecodeisalreadydeployed/fixture-nest",
		"main",
		"62139be31792ab4a43c00eadcc8af6cadd90ee66",
	)

	assert.Equal(t, "14bc77fc515e6d66b8d9c15126ee49ca55faf879", *changeString)
	assert.Equal(t, 723*time.Hour+39*time.Minute+44*time.Second+500*time.Millisecond, duration)
}

func TestObserveGitSources(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	datastore.ExpectVersionQuery(mock)

	gdb, err := datastore.OpenGormDB(db)
	assert.Nil(t, err)

	url := "https://github.com/thecodeisalreadydeployed/fixture-nest"
	branch := "main"

	// 1 commit before main
	hash := "62139be31792ab4a43c00eadcc8af6cadd90ee66"
	msg := "feat: init NestJS project"
	author := "trif0lium"

	// main commit
	revisedHash := "14bc77fc515e6d66b8d9c15126ee49ca55faf879"
	revisedMsg := "chore(app): Hello World -> fixture-nest"
	revisedAuthor := "trif0lium"

	// Return fresh app.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE observable = ?")).
		WithArgs(true).
		WillReturnRows(getObservableAppRows(hash, msg, author, url, branch))

	// Return observable of same fresh app.
	rows := sqlmock.NewRows([]string{"Observable"}).AddRow(true)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Observable FROM `apps` WHERE `apps`.`id` = ?")).
		WithArgs("app-test").
		WillReturnRows(rows)

	// Return saved app.
	expectSaveApp(mock, revisedHash, revisedMsg, revisedAuthor, url, branch)

	mock.ExpectClose()

	logger := zaptest.NewLogger(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workloadController := mock_workloadcontroller.NewMockWorkloadController(ctrl)
	workloadController.EXPECT().NewDeployment(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	repositoryObserver := NewRepositoryObserver(logger, gdb, workloadController)
	go repositoryObserver.ObserveGitSources()
	timeout := time.After(2 * time.Second)
	<-timeout

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func expectSaveApp(mock sqlmock.Sqlmock, revisedHash string, revisedMsg string, revisedAuthor string, url string, branch string) {
	query := "SELECT * FROM `apps` WHERE id = ? ORDER BY `apps`.`id` LIMIT 1"
	exec := "UPDATE `apps` SET `project_id`=?,`name`=?,`git_source`=?,`created_at`=?,`updated_at`=?,`build_configuration`=?,`observable`=? WHERE `id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs(
			"prj-test",
			"Fixture Nest",
			model.GetGitSourceString(model.GitSource{
				CommitSHA:        revisedHash,
				CommitMessage:    revisedMsg,
				CommitAuthorName: revisedAuthor,
				RepositoryURL:    url,
				Branch:           branch,
			}),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			model.GetBuildConfigurationString(model.BuildConfiguration{}),
			true,
			"app-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app-test").
		WillReturnRows(getObservableAppRows(revisedHash, revisedMsg, revisedAuthor, url, branch))
}

func getObservableAppRows(hash string, msg string, author string, url string, branch string) *sqlmock.Rows {
	app := model.App{
		ID:        "app-test",
		ProjectID: "prj-test",
		Name:      "Fixture Nest",
		GitSource: model.GitSource{
			CommitSHA:        hash,
			CommitMessage:    msg,
			CommitAuthorName: author,
			RepositoryURL:    url,
			Branch:           branch,
		},
		BuildConfiguration: model.BuildConfiguration{},
		Observable:         true,
	}

	a := datamodel.NewAppFromModel(&app)
	return sqlmock.NewRows(datastore.AppStructString()).AddRow(
		a.ID,
		a.ProjectID,
		a.Name,
		a.GitSource,
		a.CreatedAt,
		a.UpdatedAt,
		a.BuildConfiguration,
		a.Observable,
	)
}
