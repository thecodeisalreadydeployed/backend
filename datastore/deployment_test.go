package datastore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
)

func TestGetDeploymentByAppID(t *testing.T) {
	fmt.Println(datamodel.DeploymentStructString())
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `deployments` WHERE `deployments`.`app_id` = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app-test").
		WillReturnRows(getDeploymentRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetDeploymentsByAppID(gdb, "app-test")
	assert.Nil(t, err)

	expected := &[]model.Deployment{*getExpectedDeployment()}
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetDeploymentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `deployments` WHERE id = ? ORDER BY `deployments`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("dpl-test").
		WillReturnRows(getDeploymentRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetDeploymentByID(gdb, "dpl-test")
	assert.Nil(t, err)

	expected := getExpectedDeployment()
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSetDeploymentState(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)
	fmt.Println(datamodel.DeploymentStructString())

	exec := "UPDATE `deployments` SET `state`=? WHERE `deployments`.`id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(exec)).
		WithArgs(model.DeploymentStateReady, "dpl-test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	err = SetDeploymentState(gdb, "dpl-test", model.DeploymentStateReady)
	assert.Nil(t, err)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
