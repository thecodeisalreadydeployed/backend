package datastore

import (
	"bou.ke/monkey"
	"encoding/base64"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"regexp"
	"testing"
	"time"
)

func TestGetAllApps(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps`")).WillReturnRows(getAppRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetAllApps(gdb)
	assert.Nil(t, err)

	expected := &[]model.App{*getExpectedApp()}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetObservableApps(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE observable = ?")).
		WithArgs(true).
		WillReturnRows(getAppRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetObservableApps(gdb)
	assert.Nil(t, err)

	expected := &[]model.App{*getExpectedApp()}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetAppByProjectID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `apps` WHERE `apps`.`project_id` = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("prj_test").
		WillReturnRows(getAppRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetAppsByProjectID(gdb, "prj_test")
	assert.Nil(t, err)

	expected := &[]model.App{*getExpectedApp()}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetAppByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `apps` WHERE id = ? ORDER BY `apps`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app_test").
		WillReturnRows(getAppRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetAppByID(gdb, "app_test")
	assert.Nil(t, err)

	expected := getExpectedApp()
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSaveApp(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Unix(0, 0)
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `apps` WHERE id = ? ORDER BY `apps`.`id` LIMIT 1"
	exec := "UPDATE `apps` SET `project_id`=?,`name`=?,`git_source`=?,`created_at`=?,`updated_at`=?,`build_configuration`=?,`observable`=? WHERE `id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs(
			"prj_test",
			"Best App",
			getGitSourceString(),
			time.Unix(0, 0),
			time.Unix(0, 0),
			getBuildConfigurationString(),
			true,
			"app_test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app_test").
		WillReturnRows(getAppRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	expected := getExpectedApp()

	actual, err := SaveApp(gdb, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestRemoveApp(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Unix(0, 0)
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `apps` WHERE `apps`.`id` = ? ORDER BY `apps`.`id` LIMIT 1"
	exec := "DELETE FROM `apps` WHERE `apps`.`id` = ?"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app_test").
		WillReturnRows(getAppRows())

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs("app_test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	err = RemoveApp(gdb, "app_test")
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetAppByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE `apps`.`name` = ?")).
		WithArgs("Best App").
		WillReturnRows(getAppRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetAppsByName(gdb, "Best App")
	assert.Nil(t, err)

	expected := &[]model.App{*getExpectedApp()}

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func getGitSourceString() string {
	gitSource, err := json.Marshal(model.GitSource{})
	if err != nil {
		panic(err)
	}
	return cast.ToString(gitSource)
}

func getBuildConfigurationString() string {
	buildConfiguration, err := json.Marshal(model.BuildConfiguration{})
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(buildConfiguration)
}

func getAppRows() *sqlmock.Rows {
	return sqlmock.NewRows(datamodel.AppStructString()).
		AddRow(
			"app_test",
			"prj_test",
			"Best App",
			getGitSourceString(),
			time.Unix(0, 0),
			time.Unix(0, 0),
			getBuildConfigurationString(),
			true,
		)
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
