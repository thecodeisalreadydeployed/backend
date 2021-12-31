package datastore

import (
	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/model"
	"regexp"
	"testing"
	"time"
)

func TestGetEventsByDeploymentID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `events` WHERE `events`.`deployment_id` = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("dpl-test").
		WillReturnRows(GetEventRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	actual, err := GetEventsByDeploymentID(gdb, "dpl-test")
	assert.Nil(t, err)

	expected := &[]model.Event{*GetExpectedEvent()}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetEventByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `events` WHERE id = ? ORDER BY `events`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("abcdefghijklmnopqrstuvwxyz0").
		WillReturnRows(GetEventRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	actual, err := GetEventByID(gdb, "abcdefghijklmnopqrstuvwxyz0")
	assert.Nil(t, err)

	expected := GetExpectedEvent()
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSaveEvent(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Unix(0, 0)
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `events` WHERE id = ? ORDER BY `events`.`id` LIMIT 1"
	exec := "UPDATE `events` SET `deployment_id`=?,`text`=?,`type`=?,`exported_at`=?,`created_at`=? WHERE `id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs(
			"dpl-test",
			"Downloading dependencies (1/20)",
			model.INFO,
			time.Unix(0, 0),
			time.Unix(0, 0),
			"abcdefghijklmnopqrstuvwxyz0").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("abcdefghijklmnopqrstuvwxyz0").
		WillReturnRows(GetEventRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	expected := GetExpectedEvent()

	actual, err := SaveEvent(gdb, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
