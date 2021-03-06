package datastore

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thecodeisalreadydeployed/model"
	"regexp"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllProjects(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects`")).
		WillReturnRows(GetProjectRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetAllProjects()
	assert.Nil(t, err)

	expected := &[]model.Project{*GetExpectedProject()}

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetProjectByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `projects` WHERE id = ? ORDER BY `projects`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("prj-test").
		WillReturnRows(GetProjectRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetProjectByID("prj-test")
	assert.Nil(t, err)

	expected := GetExpectedProject()

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSaveProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `projects` WHERE id = ? ORDER BY `projects`.`id` LIMIT 1"
	exec := "UPDATE `projects` SET `name`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs("Best Project", sqlmock.AnyArg(), sqlmock.AnyArg(), "prj-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("prj-test").
		WillReturnRows(GetProjectRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	expected := GetExpectedProject()

	actual, err := d.SaveProject(expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestRemoveProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("prj-test").
		WillReturnRows(GetProjectRows())

	exec := "DELETE FROM `projects` WHERE `projects`.`id` = ?"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs("prj-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	err = d.RemoveProject("prj-test")
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetProjectByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE name LIKE ?")).
		WithArgs("%Best Project%").
		WillReturnRows(GetProjectRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetProjectsByName("Best Project")
	assert.Nil(t, err)

	expected := &[]model.Project{*GetExpectedProject()}

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
