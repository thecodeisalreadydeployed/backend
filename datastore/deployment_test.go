package datastore

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"regexp"
	"testing"
)

func TestGetDeploymentByAppID(t *testing.T) {
	fmt.Println(datamodel.DeploymentStructString())
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	expectVersionQuery(mock)

	query := "SELECT * FROM `deployments` WHERE `deployments`.`app_id` = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("app_test").
		WillReturnRows(getDeploymentRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetDeploymentsByAppID(gdb, "app_test")
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
		WithArgs("dpl_test").
		WillReturnRows(getDeploymentRows())
	mock.ExpectClose()

	gdb, err := openGormDB(db)
	assert.Nil(t, err)

	actual, err := GetDeploymentByID(gdb, "dpl_test")
	assert.Nil(t, err)

	expected := getExpectedDeployment()
	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
