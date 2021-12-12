package datastore

import (
	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"regexp"
	"time"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllProjects(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects`")).
		WillReturnRows(getProjectRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetAllProjects(gdb)
	assert.Nil(t, err)

	expected := &[]model.Project{*getExpectedProject()}

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetProjectByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `projects` WHERE id = ? ORDER BY `projects`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("prj_test").
		WillReturnRows(getProjectRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetProjectByID(gdb, "prj_test")
	assert.Nil(t, err)

	expected := getExpectedProject()

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSaveProject(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Unix(0, 0)
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `projects` WHERE id = ? ORDER BY `projects`.`id` LIMIT 1"
	exec := "UPDATE `projects` SET `name`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs("Best Project", time.Unix(0, 0), time.Unix(0, 0), "prj_test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("prj_test").
		WillReturnRows(getProjectRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	expected := getExpectedProject()

	actual, err := SaveProject(gdb, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestRemoveProject(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Unix(0, 0)
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1"
	exec := "DELETE FROM `projects` WHERE `projects`.`id` = ?"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("prj_test").
		WillReturnRows(getProjectRows())
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs("prj_test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectClose()

	gdb, err := openGormDB(db)

	err = RemoveProject(gdb, "prj_test")
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetProjectByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`name` = ?")).
		WithArgs("Best Project").
		WillReturnRows(getProjectRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)

	actual, err := GetProjectsByName(gdb, "Best Project")
	assert.Nil(t, err)

	expected := &[]model.Project{*getExpectedProject()}

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func getProjectRows() *sqlmock.Rows {
	return sqlmock.NewRows(datamodel.ProjectStructString()).
		AddRow("prj_test", "Best Project", time.Unix(0, 0), time.Unix(0, 0))
}

func getExpectedProject() *model.Project {
	return &model.Project{
		ID:        "prj_test",
		Name:      "Best Project",
		CreatedAt: time.Unix(0, 0),
		UpdatedAt: time.Unix(0, 0),
	}
}
