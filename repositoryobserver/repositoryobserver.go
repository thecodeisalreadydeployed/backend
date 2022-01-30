package repositoryobserver

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"
)

type RepositoryObserver interface {
	ObserveGitSources()
}

type repositoryObserver struct {
	db                 *gorm.DB
	workloadController *workloadcontroller.WorkloadController
	appChan            chan *model.App
}

func NewRepositoryObserver(DB *gorm.DB, workloadController *workloadcontroller.WorkloadController) RepositoryObserver {
	appChan := make(chan *model.App)
	return &repositoryObserver{db: DB, workloadController: workloadController, appChan: appChan}
}

func (observer *repositoryObserver) ObserveGitSources() {}

const WaitAfterErrorInterval = 10 * time.Second

func ObserveGitSources(DB *gorm.DB, observables *sync.Map, appChan chan *model.App, deploy func(string, *string) (*model.Deployment, error)) {
	for {
		apps, err := datastore.GetObservableApps(DB)
		if err != nil {
			zap.L().Error("An error occurred while accessing the database for observable apps, waiting for the next retry.\n" + err.Error())
			time.Sleep(WaitAfterErrorInterval)
		} else {
			for _, app := range *apps {
				if _, ok := observables.Load(app.ID); !ok {
					observables.Store(app.ID, nil)
					go checkGitSourceWrapper(DB, &app, observables, deploy)
				}
			}
			break
		}
	}

	for {
		app := <-appChan
		if _, ok := observables.Load(app.ID); !ok {
			observables.Store(app.ID, nil)
			go checkGitSourceWrapper(DB, app, observables, deploy)
		}
	}
}

func checkGitSourceWrapper(DB *gorm.DB, app *model.App, observables *sync.Map, deploy func(string, *string) (*model.Deployment, error)) {
	for {
		cont := checkGitSource(DB, app, observables, deploy)
		if !cont {
			return
		}
	}
}

func checkGitSource(DB *gorm.DB, app *model.App, observables *sync.Map, deploy func(string, *string) (*model.Deployment, error)) bool {
	var retry bool
	var exit bool
	for {
		retry, exit = checkObservable(DB, app, observables)
		if !retry {
			break
		}
	}
	if exit {
		return false
	}

	var commit *string
	var duration time.Duration
	var restart bool
	for {
		commit, duration = checkChanges(app.GitSource.RepositoryURL, app.GitSource.Branch, app.GitSource.CommitSHA)

		if duration > gitgateway.MaximumInterval {
			duration = gitgateway.MaximumInterval
		}
		if commit == nil {
			if duration == -1 {
				zap.L().Error(app.ID + " An error occurred while fetching the repository, waiting for the next retry.")
				time.Sleep(WaitAfterErrorInterval)
			} else {
				zap.L().Info(app.ID + " There are no changes in the application, waiting for the next repository check.")
				time.Sleep(duration)
				restart = true
				break
			}
		} else {
			restart = false
			break
		}
	}
	if restart {
		return true
	}

	for {
		_, err := deploy(app.ID, commit)
		if err != nil {
			zap.L().Error(app.ID + " An error occurred while deploying new revision of %s, waiting for the next retry.\n" + err.Error())
			time.Sleep(WaitAfterErrorInterval)
		} else {
			break
		}
	}

	zap.L().Info(app.ID + " Deployment of new revision completed, waiting for new changes.")
	time.Sleep(duration)
	return true
}

func checkObservable(DB *gorm.DB, app *model.App, observables *sync.Map) (bool, bool) {
	observableNow, err := datastore.IsObservableApp(DB, app.ID)
	if err != nil {
		zap.L().Error(app.ID + " An error occurred while accessing the database, waiting for the next retry.\n" + err.Error())
		time.Sleep(WaitAfterErrorInterval)
		return true, false
	}
	if observableNow {
		return false, false
	} else {
		observables.Delete(app.ID)
		return false, true
	}
}

func checkChanges(repoURL string, branch string, currentCommitSHA string) (*string, time.Duration) {
	git, err := gitgateway.NewGitGatewayRemote(repoURL)
	if err != nil {
		return nil, -1
	}

	duration, err := git.CommitInterval()
	if err != nil {
		return nil, -1
	}

	checkoutErr := git.Checkout(branch)
	if checkoutErr != nil {
		return nil, -1
	}

	ref, err := git.Head()
	if err != nil {
		return nil, -1
	}

	diff, diffErr := git.Diff(currentCommitSHA, ref)
	if diffErr != nil {
		return nil, -1
	}

	if len(diff) > 0 {
		return &ref, duration
	} else {
		return nil, duration
	}
}
