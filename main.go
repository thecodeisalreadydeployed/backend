package main

import (
	"github.com/thecodeisalreadydeployed/apiserver"
	"github.com/thecodeisalreadydeployed/datastore"
)

func main() {
	datastore.Init()
	apiserver.APIServer(3000)
}
