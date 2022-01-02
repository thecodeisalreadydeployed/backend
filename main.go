package main

import (
	"github.com/subosito/gotenv"
	"github.com/thecodeisalreadydeployed/apiserver"
	"github.com/thecodeisalreadydeployed/config"
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

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = []string{"stdout"}
	logger, err := loggerConfig.Build()

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)

	config.SetDefault()
	datastore.Init()
	go repositoryobserver.ObserveGitSources(datastore.GetDB())
	apiserver.APIServer(3000)
}
