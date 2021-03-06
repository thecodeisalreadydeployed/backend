package datastore

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/model"
	"regexp"
	"testing"
)

func TestGetAllApps(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps`")).WillReturnRows(GetAppRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetAllApps()
	assert.Nil(t, err)

	expected := &[]model.App{*GetExpectedApp()}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetObservableApps(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE observable = ?")).
		WithArgs(true).
		WillReturnRows(GetAppRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetObservableApps()
	assert.Nil(t, err)

	expected := &[]model.App{*GetExpectedApp()}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetAppByProjectID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `apps` WHERE `apps`.`project_id` = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("prj-test").
		WillReturnRows(GetAppRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetAppsByProjectID("prj-test")
	assert.Nil(t, err)

	expected := &[]model.App{*GetExpectedApp()}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetAppByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `apps` WHERE id = ? ORDER BY `apps`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app-test").
		WillReturnRows(GetAppRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetAppByID("app-test")
	assert.Nil(t, err)

	expected := GetExpectedApp()
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSaveApp(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `apps` WHERE id = ? ORDER BY `apps`.`id` LIMIT 1"
	exec := "UPDATE `apps` SET `project_id`=?,`name`=?,`git_source`=?,`created_at`=?,`updated_at`=?,`build_configuration`=?,`observable`=?,`fetch_interval`=? WHERE `id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs(
			"prj-test",
			"Best App",
			model.GetGitSourceString(model.GitSource{
				CommitSHA:        "a",
				CommitMessage:    "a",
				CommitAuthorName: "a",
				RepositoryURL:    "a",
				Branch:           "a",
			}),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			model.GetBuildConfigurationString(model.BuildConfiguration{}),
			true,
			0,
			"app-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app-test").
		WillReturnRows(GetAppRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	expected := GetExpectedApp()

	actual, err := d.SaveApp(expected)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestRemoveApp(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `apps` WHERE `apps`.`id` = ? ORDER BY `apps`.`id` LIMIT 1"
	exec := "DELETE FROM `apps` WHERE `apps`.`id` = ?"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app-test").
		WillReturnRows(GetAppRows())

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs("app-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	err = d.RemoveApp("app-test")
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetAppByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE name LIKE ?")).
		WithArgs("%Best App%").
		WillReturnRows(GetAppRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetAppsByName("Best App")
	assert.Nil(t, err)

	expected := &[]model.App{*GetExpectedApp()}

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestIsObservableApp(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	rows := sqlmock.NewRows([]string{"Observable"}).AddRow(true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Observable FROM `apps` WHERE `apps`.`id` = ?")).
		WithArgs("app-test").
		WillReturnRows(rows)
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.IsObservableApp("app-test")
	assert.Nil(t, err)

	assert.Equal(t, true, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSetObservable(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `apps` WHERE id = ? ORDER BY `apps`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app-test").
		WillReturnRows(GetAppRows())

	exec := "UPDATE `apps` SET `project_id`=?,`name`=?,`git_source`=?,`created_at`=?,`updated_at`=?,`build_configuration`=?,`observable`=?,`fetch_interval`=? WHERE `id` = ?"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs(
			"prj-test",
			"Best App",
			model.GetGitSourceString(model.GitSource{
				CommitSHA:        "a",
				CommitMessage:    "a",
				CommitAuthorName: "a",
				RepositoryURL:    "a",
				Branch:           "a",
			}),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			model.GetBuildConfigurationString(model.BuildConfiguration{}),
			false,
			0,
			"app-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	err = d.SetObservable("app-test", false)
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
