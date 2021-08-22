package main

import (
	"github.com/thecodeisalreadydeployed/apiserver"
	"github.com/thecodeisalreadydeployed/datastore"
)

func main() {
	db := datastore.GetDB()
	datastore.InitDB(db)
	//p := datastore.GetProject(db, &datamodel.Project{ID: "11"})
	//fmt.Println(p)
	apiserver.APIServer(3000)
}
