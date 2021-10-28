package datastore

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thecodeisalreadydeployed/datamodel"
	"github.com/thecodeisalreadydeployed/model"
	"regexp"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGetAllProjects(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	rows := sqlmock.NewRows(datamodel.ProjectStructString()).
		AddRow("prj_test", "Best Project", time.Unix(0, 0), time.Unix(0, 0))

	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mock.NewRows([]string{"version"}).FromCSVString("1"),
	)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects`")).WillReturnRows(rows)
	mock.ExpectClose()

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	actual, err := GetAllProjects(gdb)
	assert.Nil(t, err)

	expected := &[]model.Project{{
		ID:        "prj_test",
		Name:      "Best Project",
		CreatedAt: time.Unix(0, 0),
		UpdatedAt: time.Unix(0, 0),
	}}

	assert.Equal(t, expected, actual)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
