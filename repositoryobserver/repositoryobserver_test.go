package repositoryobserver

import (
	"bou.ke/monkey"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
	mock_workloadcontroller "github.com/thecodeisalreadydeployed/workloadcontroller/v2/mock"
	"go.uber.org/zap/zaptest"

	"github.com/stretchr/testify/assert"
)

func TestCheckChanges(t *testing.T) {
	observer := NewRepositoryObserver(nil, nil, nil)
	changeString, duration := observer.CheckChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		"main",
		"37e8e4d20d889924780f2373453a246591b6b11a",
	)

	assert.Equal(t, "5da29979c5ef986dc8ec6aa603e0862310abc96e", *changeString)
	assert.Equal(t, 19*time.Minute+57*time.Second, duration)

	changeString, duration = observer.CheckChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		"main",
		"5da29979c5ef986dc8ec6aa603e0862310abc96e",
	)

	assert.Nil(t, changeString)
	assert.Equal(t, 19*time.Minute+57*time.Second, duration)

	changeString, duration = observer.CheckChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-nest",
		"main",
		"62139be31792ab4a43c00eadcc8af6cadd90ee66",
	)

	assert.Equal(t, "14bc77fc515e6d66b8d9c15126ee49ca55faf879", *changeString)
	assert.Equal(t, 723*time.Hour+39*time.Minute+44*time.Second+500*time.Millisecond, duration)

	changeString, duration = observer.CheckChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-nest",
		"dev",
		"62139be31792ab4a43c00eadcc8af6cadd90ee66",
	)

	assert.Equal(t, "14bc77fc515e6d66b8d9c15126ee49ca55faf879", *changeString)
	assert.Equal(t, 723*time.Hour+39*time.Minute+44*time.Second+500*time.Millisecond, duration)
}

func TestObserveGitSources(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Unix(0, 0)
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	datastore.ExpectVersionQuery(mock)

	gdb, err := datastore.OpenGormDB(db)
	assert.Nil(t, err)

	clean, path, msg, hash, revisedMsg, revisedHash := initRepository(t)
	defer clean()

	// Return fresh app.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE observable = ?")).
		WithArgs(true).
		WillReturnRows(getObservableAppRows(path, msg, hash))

	// Return observable of same fresh app.
	rows := sqlmock.NewRows([]string{"Observable"}).AddRow(true)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Observable FROM `apps` WHERE `apps`.`id` = ?")).
		WithArgs("app-test").
		WillReturnRows(rows)

	// Return saved app.
	expectSaveApp(mock, revisedHash, revisedMsg, path)

	// Return error fetching datastore after deploying once.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Observable FROM `apps` WHERE `apps`.`id` = ?")).
		WithArgs("app-test").
		WillReturnError(errors.New("simulated failure"))

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

func expectSaveApp(mock sqlmock.Sqlmock, revisedHash string, revisedMsg string, path string) {
	query := "SELECT * FROM `apps` WHERE id = ? ORDER BY `apps`.`id` LIMIT 1"
	exec := "UPDATE `apps` SET `project_id`=?,`name`=?,`git_source`=?,`created_at`=?,`updated_at`=?,`build_configuration`=?,`observable`=? WHERE `id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs(
			"prj-test",
			"Best App",
			model.GetGitSourceString(model.GitSource{
				CommitSHA:        revisedHash,
				CommitMessage:    revisedMsg,
				CommitAuthorName: config.DefaultGitSignature().Name,
				RepositoryURL:    path,
				Branch:           "main",
			}),
			time.Unix(0, 0),
			time.Unix(0, 0),
			model.GetBuildConfigurationString(model.BuildConfiguration{}),
			true,
			"app-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app-test").
		WillReturnRows(getObservableAppRows(path, revisedMsg, revisedHash))
}

func getObservableAppRows(path string, msg string, hash string) *sqlmock.Rows {
	app := model.App{
		ID:        "app-test",
		ProjectID: "prj-test",
		Name:      "Best App",
		GitSource: model.GitSource{
			CommitSHA:        hash,
			CommitMessage:    msg,
			CommitAuthorName: config.DefaultGitSignature().Name,
			RepositoryURL:    path,
			Branch:           "main",
		},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
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

func initRepository(t *testing.T) (func(), string, string, string, string, string) {
	path, clean := gitgateway.InitRepository()

	git, err := gitgateway.NewGitGatewayLocal(path)
	assert.Nil(t, err)

	err = git.WriteFile(".thecodeisalreadydeployed", "data")
	assert.Nil(t, err)

	msg := "This is a commit."
	hash, err := git.Commit([]string{".thecodeisalreadydeployed"}, msg)
	assert.Nil(t, err)

	err = git.WriteFile(".thecodeisalreadydeployed", "new data")
	assert.Nil(t, err)

	revisedMsg := "This is another commit."
	revisedHash, err := git.Commit([]string{".thecodeisalreadydeployed"}, revisedMsg)
	assert.Nil(t, err)

	return clean, path, msg, hash, revisedMsg, revisedHash
}
