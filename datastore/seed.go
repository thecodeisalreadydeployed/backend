package datastore

import (
	"fmt"
	faker "github.com/bxcodec/faker/v3"
	"github.com/thecodeisalreadydeployed/datamodel"
	"gorm.io/gorm"
)

func seed(db *gorm.DB) {
	seedProjects(db, 20)
}

func seedProjects(db *gorm.DB, size int) {
	exists := db.Table("projects").Find(&datamodel.Project{})
	if exists == nil {
		var ps []datamodel.Project
		for i := 0; i < size; i++ {
			var p datamodel.Project
			err := faker.FakeData(&p)
			if err != nil {
				fmt.Println(err)
			}
			ps = append(ps, p)
		}
		db.Create(&ps)
	} else {
		fmt.Println("Table Project already seeded.")
	}
}

func seedApps(db *gorm.DB) {

}

func seedDeployments(db *gorm.DB) {

}

