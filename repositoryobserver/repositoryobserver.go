package repositoryobserver

import (
	"context"
	"github.com/spf13/cast"
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
	Refresh(id string)
}

type repositoryObserver struct {
	logger             *zap.Logger
	db                 *gorm.DB
	workloadController workloadcontroller.WorkloadController
	appChan            chan *model.App
	refreshChan        map[string]chan bool
	refreshLock        *sync.Mutex
	refreshCtx         context.Context

	// {app.ID: true} if idle, {app.ID: false} if not idle, non-existent if not observable
	idleObservables *sync.Map
}

func NewRepositoryObserver(logger *zap.Logger, DB *gorm.DB, workloadController workloadcontroller.WorkloadController) RepositoryObserver {
	appChan := make(chan *model.App)
	refreshChan := make(map[string]chan bool)
	return &repositoryObserver{
		logger:             logger,
		db:                 DB,
		workloadController: workloadController,
		appChan:            appChan,
		refreshChan:        refreshChan,
		refreshCtx:         context.Background(),
		idleObservables:    &sync.Map{},
	}
}

const WaitAfterErrorInterval = 10 * time.Second

func (observer *repositoryObserver) ObserveGitSources() {
	for {
		apps, err := datastore.GetObservableApps(observer.db)
		if err != nil {
			observer.logger.Error("cannot get observable apps", zap.Error(err))
			time.Sleep(WaitAfterErrorInterval)
		} else {
			for _, app := range *apps {
				if _, ok := observer.idleObservables.Load(app.ID); !ok {
					observer.idleObservables.Store(app.ID, false)
					observer.refreshChan[app.ID] = make(chan bool)
					go observer.checkGitSourceWrapper(&app)
				}
			}
			break
		}
	}

	for {
		app := <-observer.appChan
		if _, ok := observer.idleObservables.Load(app.ID); !ok {
			observer.idleObservables.Store(app.ID, false)
			observer.refreshChan[app.ID] = make(chan bool)
			go observer.checkGitSourceWrapper(app)
		}
	}
}

func (observer *repositoryObserver) Refresh(id string) {
	observer.refreshLock.Lock()
	idle, ok := observer.idleObservables.Load(id)
	if ok && cast.ToBool(idle) {
		observer.refreshChan[id] <- true
	}
	observer.refreshLock.Unlock()
}

func (observer *repositoryObserver) checkGitSourceWrapper(app *model.App) {
	for {
		cont := observer.checkGitSource(app)
		if !cont {
			return
		}
	}
}

func (observer *repositoryObserver) checkGitSource(app *model.App) bool {
	logger := observer.logger.With(zap.String("appID", app.ID))
	var retry bool
	var exit bool
	for {
		retry, exit = observer.checkObservable(logger, app)
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
				logger.Error("failed to fetch the repository, waiting for the next retry")
				time.Sleep(WaitAfterErrorInterval)
			} else {
				logger.Info("no changes in the application, waiting for the next repository check")
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
		_, err := observer.workloadController.NewDeployment(app.ID, commit)
		if err != nil {
			logger.Error("failed to deploy new revision, waiting for the next retry", zap.Error(err))
			time.Sleep(WaitAfterErrorInterval)
		} else {
			break
		}
	}

	observer.idleObservables.Store(app.ID, true)
	select {
	case <-time.After(duration):
		break
	case <-observer.refreshChan[app.ID]:
		break
	}
	observer.refreshLock.
		observer.idleObservables.Store(app.ID, false)
	observer.refreshLock.Unlock()
	return true
}

func (observer *repositoryObserver) checkObservable(logger *zap.Logger, app *model.App) (bool, bool) {
	observableNow, err := datastore.IsObservableApp(observer.db, app.ID)
	if err != nil {
		logger.Error("application status check failed", zap.Error(err))
		time.Sleep(WaitAfterErrorInterval)
		return true, false
	}
	if observableNow {
		return false, false
	} else {
		observer.idleObservables.Delete(app.ID)
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
