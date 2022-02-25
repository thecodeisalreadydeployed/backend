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
	Refresh(id string)
	CheckChanges(repoURL string, branch string, currentCommitSHA string) (*string, time.Duration)
}

type repositoryObserver struct {
	logger             *zap.Logger
	db                 *gorm.DB
	workloadController workloadcontroller.WorkloadController
	refreshChan        map[string]chan bool
	observables        *sync.Map
}

func NewRepositoryObserver(logger *zap.Logger, DB *gorm.DB, workloadController workloadcontroller.WorkloadController) RepositoryObserver {
	refreshChan := make(map[string]chan bool)
	return &repositoryObserver{
		logger:             logger,
		db:                 DB,
		workloadController: workloadController,
		refreshChan:        refreshChan,
		observables:        &sync.Map{},
	}
}

const waitAfterErrorInterval = 10 * time.Second
const sleepObserverInterval = 30 * time.Second

func (observer *repositoryObserver) ObserveGitSources() {
	for {
		apps, err := datastore.GetObservableApps(observer.db)
		if err != nil {
			observer.logger.Error("cannot get observable apps", zap.Error(err))
			time.Sleep(waitAfterErrorInterval)
		} else {
			observer.logger.Info("obtained observable apps")
			for _, app := range *apps {
				if _, ok := observer.observables.Load(app.ID); !ok {
					observer.observables.Store(app.ID, nil)
					observer.refreshChan[app.ID] = make(chan bool)
					observer.checkGitSourceWrapper(&app)
				}
			}
		}
		time.Sleep(sleepObserverInterval)
	}
}

func (observer *repositoryObserver) Refresh(id string) {
	_, ok := observer.observables.Load(id)
	if ok {
		select {
		case observer.refreshChan[id] <- true:
			break
		default:
			break
		}
	}
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
		commit, duration = observer.CheckChanges(app.GitSource.RepositoryURL, app.GitSource.Branch, app.GitSource.CommitSHA)

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
	if restart {
		return true
	}

	for {
		_, err := observer.workloadController.NewDeployment(app.ID, commit)
		if err != nil {
			logger.Error("failed to deploy new revision, waiting for the next retry", zap.Error(err))
			time.Sleep(waitAfterErrorInterval)
		} else {
			break
		}
	}
	logger.Info("deployment completed")

	select {
	case <-observer.refreshChan[app.ID]:
		break
	case <-time.After(duration):
		break
	}
	return true
}

func (observer *repositoryObserver) checkObservable(logger *zap.Logger, app *model.App) (bool, bool) {
	observableNow, err := datastore.IsObservableApp(observer.db, app.ID)
	if err != nil {
		logger.Error("application status check failed", zap.Error(err))
		time.Sleep(waitAfterErrorInterval)
		return true, false
	}
	if observableNow {
		return false, false
	} else {
		observer.observables.Delete(app.ID)
		return false, true
	}
}

func (observer *repositoryObserver) CheckChanges(repoURL string, branch string, currentCommitSHA string) (*string, time.Duration) {
	git, err := gitgateway.NewGitGatewayRemote(repoURL)
	if err != nil {
		observer.logger.Error("cannot connect to remote", zap.Error(err))
		return nil, -1
	}

	duration, err := git.CommitInterval()
	if err != nil {
		observer.logger.Error("cannot get commit interval", zap.Error(err))
		return nil, -1
	}

	checkoutErr := git.Checkout(branch)
	if checkoutErr != nil {
		observer.logger.Error("cannot checkout", zap.Error(checkoutErr))
		return nil, -1
	}

	ref, err := git.Head()
	if err != nil {
		observer.logger.Error("cannot get repository head", zap.Error(err))
		return nil, -1
	}

	diff, diffErr := git.Diff(currentCommitSHA, ref)
	if diffErr != nil {
		observer.logger.Error("cannot get commit diff", zap.Error(diffErr))
		return nil, -1
	}

	if len(diff) > 0 {
		return &ref, duration
	} else {
		return nil, duration
	}
}
