package repositoryobserver

import (
	"github.com/thecodeisalreadydeployed/gitapi"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
	"github.com/thecodeisalreadydeployed/workloadcontroller/v2"
)

type RepositoryObserver interface {
	ObserveGitSources()
	Refresh(id string) bool
	CheckChanges(logger *zap.Logger, repoURL string, branch string, currentCommitSHA string) (*string, time.Duration)
}

type repositoryObserver struct {
	logger             *zap.Logger
	dataStore          datastore.DataStore
	workloadController workloadcontroller.WorkloadController
	refreshChan        map[string]chan bool
	observables        *sync.Map
	gapi               gitapi.GitAPIBackend
}

func NewRepositoryObserver(logger *zap.Logger, dataStore datastore.DataStore, workloadController workloadcontroller.WorkloadController) RepositoryObserver {
	refreshChan := make(map[string]chan bool)
	gapi := gitapi.NewGitAPIBackend(logger)
	return &repositoryObserver{
		logger:             logger,
		dataStore:          dataStore,
		workloadController: workloadController,
		refreshChan:        refreshChan,
		observables:        &sync.Map{},
		gapi:               gapi,
	}
}

const waitAfterErrorInterval = 10 * time.Second
const sleepObserverInterval = 30 * time.Second

func (observer *repositoryObserver) ObserveGitSources() {
	for {
		apps, err := observer.dataStore.GetObservableApps()
		if err != nil {
			observer.logger.Error("cannot get observable apps", zap.Error(err))
			time.Sleep(waitAfterErrorInterval)
		} else {
			for _, app := range *apps {
				if _, ok := observer.observables.Load(app.ID); !ok {
					thisApp := app
					observer.observables.Store(thisApp.ID, nil)
					observer.refreshChan[thisApp.ID] = make(chan bool)
					go observer.checkGitSourceWrapper(&thisApp)
				}
			}
		}
		time.Sleep(sleepObserverInterval)
	}
}

func (observer *repositoryObserver) Refresh(id string) bool {
	_, ok := observer.observables.Load(id)
	if ok {
		select {
		case observer.refreshChan[id] <- true:
			return true
		default:
			return false
		}
	}
	return false
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

	commit, duration, restart := observer.reportChanges(app, logger)
	if restart {
		return true
	}

	observer.newDeployment(logger, app, commit)

	observer.fillGitSource(logger, app)

	observer.saveApp(logger, app)

	logger.Info("deployment completed")

	select {
	case <-observer.refreshChan[app.ID]:
		logger.Info("refreshing: now observing for changes")
		break
	case <-time.After(duration):
		break
	}
	return true
}

func (observer *repositoryObserver) reportChanges(app *model.App, logger *zap.Logger) (*string, time.Duration, bool) {
	var commit *string
	var duration time.Duration
	var restart bool
	for {
		commit, duration = observer.CheckChanges(logger, app.GitSource.RepositoryURL, app.GitSource.Branch, app.GitSource.CommitSHA)

		if app.FetchInterval != 0 {
			duration = time.Duration(app.FetchInterval) * time.Second
		}

		if duration > gitgateway.MaximumInterval {
			duration = gitgateway.MaximumInterval
		}
		if commit == nil {
			if duration == -1 {
				logger.Error("failed to fetch the repository, waiting for the next retry")
				time.Sleep(waitAfterErrorInterval)
			} else {
				logger.Info("no changes in the application, waiting for the next repository check")
				select {
				case <-observer.refreshChan[app.ID]:
					logger.Info("refreshing: now observing for changes")
					break
				case <-time.After(duration):
					break
				}
				restart = true
				break
			}
		} else {
			restart = false
			break
		}
	}
	return commit, duration, restart
}

func (observer *repositoryObserver) newDeployment(logger *zap.Logger, app *model.App, commit *string) {
	for {
		_, err := observer.workloadController.NewDeployment(app.ID, commit, observer.dataStore)
		if err != nil {
			logger.Error("failed to deploy new revision, waiting for the next retry", zap.Error(err))
			time.Sleep(waitAfterErrorInterval)
		} else {
			break
		}
	}
}

func (observer *repositoryObserver) fillGitSource(logger *zap.Logger, app *model.App) {
	for {
		gs, err := observer.gapi.FillGitSource(&app.GitSource)
		if err != nil {
			logger.Error("failed to get new commit info, waiting for the next retry", zap.Error(err))
			time.Sleep(waitAfterErrorInterval)
		} else {
			app.GitSource = *gs
			break
		}
	}
}

func (observer *repositoryObserver) saveApp(logger *zap.Logger, app *model.App) {
	for {
		_, err := observer.dataStore.SaveApp(app)
		if err != nil {
			logger.Error("failed to save new commit info, waiting for the next retry", zap.Error(err))
			time.Sleep(waitAfterErrorInterval)
		} else {
			break
		}
	}
}

func (observer *repositoryObserver) checkObservable(logger *zap.Logger, app *model.App) (bool, bool) {
	observableNow, err := observer.dataStore.IsObservableApp(app.ID)
	if err != nil {
		logger.Error("application status check failed", zap.Error(err))
		time.Sleep(waitAfterErrorInterval)
		return true, false
	}
	if observableNow {
		return false, false
	} else {
		logger.Info("app is now set to not be observed")
		observer.observables.Delete(app.ID)
		return false, true
	}
}

func (observer *repositoryObserver) CheckChanges(logger *zap.Logger, repoURL string, branch string, currentCommitSHA string) (*string, time.Duration) {
	git, err := gitgateway.NewGitGatewayRemote(repoURL)
	if err != nil {
		logger.Error("cannot connect to remote", zap.Error(err))
		return nil, -1
	}

	duration, err := git.CommitInterval()
	if err != nil {
		logger.Error("cannot get commit interval", zap.Error(err))
		return nil, -1
	}

	checkoutErr := git.Checkout(branch)
	if checkoutErr != nil {
		logger.Error("cannot checkout", zap.Error(checkoutErr))
		return nil, -1
	}

	ref, err := git.Head()
	if err != nil {
		logger.Error("cannot get repository head", zap.Error(err))
		return nil, -1
	}

	diff, diffErr := git.Diff(currentCommitSHA, ref)
	if diffErr != nil {
		logger.Error("cannot get commit diff", zap.Error(diffErr))
		return nil, -1
	}

	if len(diff) > 0 {
		return &ref, duration
	} else {
		return nil, duration
	}
}
