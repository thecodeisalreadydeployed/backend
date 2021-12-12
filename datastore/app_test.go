package datastore

import (
	"encoding/base64"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

func TestGetAllApps(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	gitSource, err := json.Marshal(model.GitSource{})
	if err != nil {
		panic(err)
	}

	buildConfiguration, err := json.Marshal(model.BuildConfiguration{})
	if err != nil {
		panic(err)
	}

	rows := sqlmock.NewRows(datamodel.AppStructString()).
		AddRow(
			"app_test",
			"prj_test",
			"Best App",
			cast.ToString(gitSource),
			time.Unix(0, 0),
			time.Unix(0, 0),
			base64.StdEncoding.EncodeToString(buildConfiguration),
			true,
		)

	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mock.NewRows([]string{"version"}).FromCSVString("1"),
	)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps`")).WillReturnRows(rows)
	mock.ExpectClose()

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	actual, err := GetAllApps(gdb)
	assert.Nil(t, err)

	expected := &[]model.App{{
		ID:                 "app_test",
		ProjectID:          "prj_test",
		Name:               "Best App",
		GitSource:          model.GitSource{},
		CreatedAt:          time.Unix(0, 0),
		UpdatedAt:          time.Unix(0, 0),
		BuildConfiguration: model.BuildConfiguration{},
		Observable:         true,
	}}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

//func TestGetProjectByID(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	assert.Nil(t, err)
//
//	rows := sqlmock.NewRows(datamodel.ProjectStructString()).
//		AddRow("prj_test", "Best Project", time.Unix(0, 0), time.Unix(0, 0))
//
//	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
//		mock.NewRows([]string{"version"}).FromCSVString("1"),
//	)
//
//	query := "SELECT * FROM `projects` WHERE id = ? ORDER BY `projects`.`id` LIMIT 1"
//	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("prj_test").WillReturnRows(rows)
//	mock.ExpectClose()
//
//	gdb, err := gorm.Open(mysql.New(mysql.Config{
//		Conn: db,
//	}), &gorm.Config{})
//	assert.Nil(t, err)
//
//	actual, err := GetProjectByID(gdb, "prj_test")
//	assert.Nil(t, err)
//
//	expected := &model.Project{
//		ID:        "prj_test",
//		Name:      "Best Project",
//		CreatedAt: time.Unix(0, 0),
//		UpdatedAt: time.Unix(0, 0),
//	}
//
//	assert.Equal(t, expected, actual)
//
//	err = db.Close()
//	assert.Nil(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.Nil(t, err)
//}
//
//func TestSaveProject(t *testing.T) {
//	monkey.Patch(time.Now, func() time.Time {
//		return time.Unix(0, 0)
//	})
//	defer monkey.UnpatchAll()
//
//	db, mock, err := sqlmock.New()
//	assert.Nil(t, err)
//
//	rows := sqlmock.NewRows(datamodel.ProjectStructString()).
//		AddRow("prj_test", "Best Project", time.Unix(0, 0), time.Unix(0, 0))
//
//	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
//		mock.NewRows([]string{"version"}).FromCSVString("1"),
//	)
//
//	query := "SELECT * FROM `projects` WHERE id = ? ORDER BY `projects`.`id` LIMIT 1"
//	exec := "UPDATE `projects` SET `name`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?"
//
//	mock.ExpectBegin()
//	mock.ExpectExec(regexp.QuoteMeta(exec)).
//		WithArgs("Best Project", time.Unix(0, 0), time.Unix(0, 0), "prj_test").
//		WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectCommit()
//	mock.ExpectQuery(regexp.QuoteMeta(query)).
//		WithArgs("prj_test").
//		WillReturnRows(rows)
//	mock.ExpectClose()
//
//	gdb, err := gorm.Open(mysql.New(mysql.Config{
//		Conn: db,
//	}), &gorm.Config{})
//	assert.Nil(t, err)
//
//	prj := &model.Project{
//		ID:        "prj_test",
//		Name:      "Best Project",
//		CreatedAt: time.Unix(0, 0),
//		UpdatedAt: time.Unix(0, 0),
//	}
//
//	actual, err := SaveProject(gdb, prj)
//	assert.Nil(t, err)
//	assert.Equal(t, prj, actual)
//
//	err = db.Close()
//	assert.Nil(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.Nil(t, err)
//}
//
//func TestRemoveProject(t *testing.T) {
//	monkey.Patch(time.Now, func() time.Time {
//		return time.Unix(0, 0)
//	})
//	defer monkey.UnpatchAll()
//
//	db, mock, err := sqlmock.New()
//	assert.Nil(t, err)
//
//	rows := sqlmock.NewRows(datamodel.ProjectStructString()).
//		AddRow("prj_test", "Best Project", time.Unix(0, 0), time.Unix(0, 0))
//
//	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
//		mock.NewRows([]string{"version"}).FromCSVString("1"),
//	)
//
//	query := "SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1"
//	exec := "DELETE FROM `projects` WHERE `projects`.`id` = ?"
//
//	mock.ExpectQuery(regexp.QuoteMeta(query)).
//		WithArgs("prj_test").
//		WillReturnRows(rows)
//
//	mock.ExpectBegin()
//	mock.ExpectExec(regexp.QuoteMeta(exec)).
//		WithArgs("prj_test").
//		WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectCommit()
//
//	mock.ExpectClose()
//
//	gdb, err := gorm.Open(mysql.New(mysql.Config{
//		Conn: db,
//	}), &gorm.Config{})
//	assert.Nil(t, err)
//
//	err = RemoveProject(gdb, "prj_test")
//	assert.Nil(t, err)
//
//	err = db.Close()
//	assert.Nil(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.Nil(t, err)
//}
//
//func TestGetProjectByName(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	assert.Nil(t, err)
//
//	rows := sqlmock.NewRows(datamodel.ProjectStructString()).
//		AddRow("prj_test", "Best Project", time.Unix(0, 0), time.Unix(0, 0))
//
//	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
//		mock.NewRows([]string{"version"}).FromCSVString("1"),
//	)
//	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`name` = ?")).
//		WithArgs("Best Project").
//		WillReturnRows(rows)
//	mock.ExpectClose()
//
//	gdb, err := gorm.Open(mysql.New(mysql.Config{
//		Conn: db,
//	}), &gorm.Config{})
//	assert.Nil(t, err)
//
//	actual, err := GetProjectsByName(gdb, "Best Project")
//	assert.Nil(t, err)
//
//	expected := &[]model.Project{{
//		ID:        "prj_test",
//		Name:      "Best Project",
//		CreatedAt: time.Unix(0, 0),
//		UpdatedAt: time.Unix(0, 0),
//	}}
//
//	assert.Equal(t, expected, actual)
//
//	err = db.Close()
//	assert.Nil(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.Nil(t, err)
//}
