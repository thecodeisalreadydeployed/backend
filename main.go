package main

import (
	"github.com/thecodeisalreadydeployed/apiserver"
	"github.com/thecodeisalreadydeployed/datastore"
)

func main() {
	datastore.InitDB()
	apiserver.APIServer(3000)
}
