package main

import (
	"github.com/subosito/gotenv"
	"github.com/thecodeisalreadydeployed/apiserver"
	"github.com/thecodeisalreadydeployed/config"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitopscontroller"
	"github.com/thecodeisalreadydeployed/gitopscontroller/argocd"
	"github.com/thecodeisalreadydeployed/repositoryobserver"
	"github.com/thecodeisalreadydeployed/util"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"
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
	datastore.Migrate()
	argoCDClient := argocd.NewArgoCDClient()
	gitOpsController := gitopscontroller.NewGitOpsController(zap.L(), argoCDClient)
	workloadController := workloadcontroller.NewWorkloadController(zap.L(), gitOpsController)
	repositoryObserver := repositoryobserver.NewRepositoryObserver(zap.L(), datastore.GetDB(), workloadController)
	go repositoryObserver.ObserveGitSources()
	apiserver.APIServer(3000, workloadController)
}
