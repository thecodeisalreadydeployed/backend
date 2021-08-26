package main

import (
	"github.com/thecodeisalreadydeployed/apiserver"
	"github.com/thecodeisalreadydeployed/datastore"
	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)

	datastore.Init()
	apiserver.APIServer(3000)
}
