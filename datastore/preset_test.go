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

func TestGetAllPresets(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `presets`")).
		WillReturnRows(GetPresetRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	actual, err := GetAllPresets(gdb)
	assert.Nil(t, err)

	expected := &[]model.Preset{*GetExpectedPreset()}

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetPresetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `presets` WHERE id = ? ORDER BY `presets`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("pst-test").
		WillReturnRows(GetPresetRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	actual, err := GetPresetByID(gdb, "pst-test")
	assert.Nil(t, err)

	expected := GetExpectedPreset()

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSavePreset(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Unix(0, 0)
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `presets` WHERE id = ? ORDER BY `presets`.`id` LIMIT 1"
	exec := "UPDATE `presets` SET `name`=?,`template`=? WHERE `id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs("My Preset", "UlVOIGVjaG8gaGVsbG8=", "pst-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("pst-test").
		WillReturnRows(GetPresetRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	expected := GetExpectedPreset()

	actual, err := SavePreset(gdb, expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestRemovePreset(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Unix(0, 0)
	})
	defer monkey.UnpatchAll()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `presets` WHERE `presets`.`id` = ? ORDER BY `presets`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("pst-test").
		WillReturnRows(GetPresetRows())

	exec := "DELETE FROM `presets` WHERE `presets`.`id` = ?"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs("pst-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	err = RemovePreset(gdb, "pst-test")
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetPresetsByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `presets` WHERE `presets`.`name` = ?")).
		WithArgs("My Preset").
		WillReturnRows(GetPresetRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	actual, err := GetPresetsByName(gdb, "My Preset")
	assert.Nil(t, err)

	expected := &[]model.Preset{*GetExpectedPreset()}

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
