package datastore

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/model"
	"regexp"
	"testing"
)

func TestGetDeploymentByAppID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `deployments` WHERE `deployments`.`app_id` = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app-test").
		WillReturnRows(GetDeploymentRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetDeploymentsByAppID("app-test")
	assert.Nil(t, err)

	expected := &[]model.Deployment{*GetExpectedDeployment()}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetDeploymentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `deployments` WHERE id = ? ORDER BY `deployments`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("dpl-test").
		WillReturnRows(GetDeploymentRows())
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	actual, err := d.GetDeploymentByID("dpl-test")
	assert.Nil(t, err)

	expected := GetExpectedDeployment()
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSetDeploymentState(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)
	fmt.Println(DeploymentStructString())

	exec := "UPDATE `deployments` SET `state`=? WHERE `deployments`.`id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs(model.DeploymentStateReady, "dpl-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	err = d.SetDeploymentState("dpl-test", model.DeploymentStateReady)
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestRemoveDeployment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ExpectVersionQuery(mock)

	query := "SELECT * FROM `deployments` WHERE `deployments`.`id` = ? ORDER BY `deployments`.`id` LIMIT 1"
	exec := "DELETE FROM `deployments` WHERE `deployments`.`id` = ?"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("dpl-test").
		WillReturnRows(GetDeploymentRows())
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs("dpl-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectClose()

	gdb, err := OpenGormDB(db)
	assert.Nil(t, err)

	d := NewMockDataStore(gdb, t)

	err = d.RemoveDeployment("dpl-test")
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
