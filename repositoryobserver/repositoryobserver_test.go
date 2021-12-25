package repositoryobserver

import (
	"bou.ke/monkey"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
	"regexp"
	"sync"
	"testing"
	"time"

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

	changeString, _ = checkChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		"main",
		"5da29979c5ef986dc8ec6aa603e0862310abc96e",
	)

	assert.Nil(t, changeString)

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

func TestFetchObservableApps(t *testing.T) {
	monkey.Patch(time.Sleep, func(d time.Duration) {
		fmt.Println("Sleep skipped.")
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	datastore.ExpectVersionQuery(mock)

	gdb, err := datastore.OpenGormDB(db)
	assert.Nil(t, err)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE observable = ?")).
		WithArgs(true).
		WillReturnRows(datastore.GetAppRows())
	mock.ExpectClose()

	aChan := make(chan *model.App)
	var wgFetch sync.WaitGroup
	wgFetch.Add(1)
	var observables sync.Map

	go fetchObservableApps(gdb, aChan, &wgFetch, &observables)

	app := *<-aChan

	assert.Equal(t, datastore.GetExpectedApp(), &app)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestCheckGitSource(t *testing.T) {
	monkey.Patch(time.Sleep, func(d time.Duration) {
		fmt.Println("Sleep skipped.")
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	datastore.ExpectVersionQuery(mock)

	gdb, err := datastore.OpenGormDB(db)
	assert.Nil(t, err)

	rows := sqlmock.NewRows([]string{"Observable"}).AddRow(false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Observable FROM `apps` WHERE `apps`.`id` = ?")).
		WithArgs("app_test").
		WillReturnRows(rows)
	mock.ExpectClose()

	path, clean := gitgateway.InitRepository()
	defer clean()

	git, err := gitgateway.NewGitGatewayLocal(path)
	assert.Nil(t, err)

	err = git.WriteFile(".thecodeisalreadydeployed", "data")
	assert.Nil(t, err)

	hash, err := git.Commit([]string{".thecodeisalreadydeployed"}, "This is a commit.")
	assert.Nil(t, err)

	app := model.App{
		ID:        "app_test",
		ProjectID: "prj_test",
		Name:      "BestApp",
		GitSource: model.GitSource{
			CommitSHA:        hash,
			CommitMessage:    "This is a commit.",
			CommitAuthorName: config.DefaultGitSignature().Name,
			RepositoryURL:    path,
			Branch:           "main",
		},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		BuildConfiguration: model.BuildConfiguration{},
		Observable:         true,
	}

	err = git.WriteFile(".thecodeisalreadydeployed", "new data")
	assert.Nil(t, err)

	_, err = git.Commit([]string{".thecodeisalreadydeployed"}, "This is another commit.")
	assert.Nil(t, err)

	cChan := make(chan bool)
	var observables sync.Map

	go checkGitSource(gdb, app, cChan, &observables)

	cont := <-cChan
	assert.False(t, cont)
}
