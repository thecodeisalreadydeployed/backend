package datastore

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/logger"
	"gorm.io/gorm"
	"reflect"
)

func createTable(db *gorm.DB, i interface{}) {
	if !db.Migrator().HasTable(i) {
		err := db.Migrator().CreateTable(i)
		if err != nil {
			logger.Error(err.Error())
		}
	} else {
		name := reflect.TypeOf(i).Elem().Name()
		logger.Info(fmt.Sprintf("Table %s already created.\n", name))
	}
}
