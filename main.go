package main

import (
	"github.com/subosito/gotenv"
	"github.com/thecodeisalreadydeployed/apiserver"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/repositoryobserver"
	"github.com/thecodeisalreadydeployed/util"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

func main() {
	if util.IsDevEnvironment() {
		err := gotenv.Load(".env.development")
		if err != nil {
			panic(err)
		}
	}

	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)

	datastore.Init()
	go repositoryobserver.ObserveGitSources(
		datastore.GetDB(),
		datastore.GetObservables(),
		datastore.GetAppChannel(),
	)
	apiserver.APIServer(3000)
}
