package repositoryobserver

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"bou.ke/monkey"
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
	changeString, duration := checkChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		"main",
		"37e8e4d20d889924780f2373453a246591b6b11a",
	)

	assert.Equal(t, "5da29979c5ef986dc8ec6aa603e0862310abc96e", *changeString)
	assert.Equal(t, 19*time.Minute+57*time.Second, duration)

	changeString, duration = checkChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		"main",
		"5da29979c5ef986dc8ec6aa603e0862310abc96e",
	)

	assert.Nil(t, changeString)
	assert.Equal(t, 19*time.Minute+57*time.Second, duration)

	changeString, duration = checkChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-nest",
		"main",
		"62139be31792ab4a43c00eadcc8af6cadd90ee66",
	)

	assert.Equal(t, "14bc77fc515e6d66b8d9c15126ee49ca55faf879", *changeString)
	assert.Equal(t, 723*time.Hour+39*time.Minute+44*time.Second+500*time.Millisecond, duration)

	changeString, duration = checkChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-nest",
		"dev",
		"62139be31792ab4a43c00eadcc8af6cadd90ee66",
	)

	assert.Equal(t, "14bc77fc515e6d66b8d9c15126ee49ca55faf879", *changeString)
	assert.Equal(t, 723*time.Hour+39*time.Minute+44*time.Second+500*time.Millisecond, duration)
}

func TestObserveGitSources(t *testing.T) {
	monkey.Patch(time.Sleep, func(d time.Duration) {
		fmt.Println("Sleep skipped.")
	})
	defer monkey.UnpatchAll()

	now := time.Now()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	datastore.ExpectVersionQuery(mock)

	gdb, err := datastore.OpenGormDB(db)
	assert.Nil(t, err)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE observable = ?")).
		WithArgs(true).
		WillReturnError(errors.New("simulated failure"))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE observable = ?")).
		WithArgs(true).
		WillReturnRows(getObservableAppRows(t, false))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Observable FROM `apps` WHERE `apps`.`id` = ?")).
		WithArgs("app-test").
		WillReturnError(errors.New("simulated failure"))

	rows := sqlmock.NewRows([]string{"Observable"}).AddRow(true)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Observable FROM `apps` WHERE `apps`.`id` = ?")).
		WithArgs("app-test").
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"Observable"}).AddRow(false)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Observable FROM `apps` WHERE `apps`.`id` = ?")).
		WithArgs("app-test").
		WillReturnRows(rows)

	mock.ExpectClose()

	logger := zaptest.NewLogger(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workloadController := mock_workloadcontroller.NewMockWorkloadController(ctrl)
	workloadController.EXPECT().NewDeployment(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	repositoryObserver := NewRepositoryObserver(logger, gdb, workloadController)
	go repositoryObserver.ObserveGitSources()

	for {
		if time.Now().After(now.Add(2 * time.Second)) {
			err = db.Close()
			assert.Nil(t, err)

			err = mock.ExpectationsWereMet()
			assert.Nil(t, err)

			return
		}
	}
}

func getObservableAppRows(t *testing.T, revised bool) *sqlmock.Rows {
	path, clean := gitgateway.InitRepository()
	defer clean()

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

	app := model.App{
		ID:        "app-test",
		ProjectID: "prj-test",
		Name:      "BestApp",
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

	if revised {
		app.GitSource.CommitSHA = revisedHash
		app.GitSource.CommitMessage = revisedMsg
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
