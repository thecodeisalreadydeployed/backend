package repositoryobserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckChanges(t *testing.T) {

	changeString, duration := checkChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		"main",
		"37e8e4d20d889924780f2373453a246591b6b11a",
	)

	assert.Equal(t, "5da29979c5ef986dc8ec6aa603e0862310abc96e", *changeString)
	assert.Equal(t, 19*time.Minute, duration)

	changeString, duration = checkChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-monorepo",
		"main",
		"5da29979c5ef986dc8ec6aa603e0862310abc96e",
	)

	assert.Nil(t, changeString)

	changeString, duration = checkChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-nest",
		"main",
		"62139be31792ab4a43c00eadcc8af6cadd90ee66",
	)

	assert.Equal(t, "14bc77fc515e6d66b8d9c15126ee49ca55faf879", *changeString)
	assert.Equal(t, 723*time.Hour+39*time.Minute, duration)

	changeString, duration = checkChanges(
		"https://github.com/thecodeisalreadydeployed/fixture-nest",
		"dev",
		"62139be31792ab4a43c00eadcc8af6cadd90ee66",
	)

	assert.Equal(t, "14bc77fc515e6d66b8d9c15126ee49ca55faf879", *changeString)
	assert.Equal(t, 723*time.Hour+39*time.Minute, duration)
}

func TestCheckObservableApps(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	datastore.ExpectVersionQuery(mock)

	gdb, err := datastore.OpenGormDB(db)
	assert.Nil(t, err)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `apps` WHERE observable = ?")).
		WithArgs(true).
		WillReturnRows(datastore.GetAppRows())
	mock.ExpectClose()

	aChan := make(chan *model.App)

	go checkObservableApps(gdb, aChan, false)

	app := *<-aChan

	assert.Equal(t, datastore.GetExpectedApp(), &app)

	err = db.Close()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
