package repositoryobserver

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
)

const MaximumDuration = 3 * time.Minute
const WaitAfterErrorInterval = 10 * time.Second

func ObserveGitSources(DB *gorm.DB, observables *sync.Map, appChan chan *model.App) {
	for {
		apps, err := datastore.GetObservableApps(DB)
		if err != nil {
			zap.L().Error("An error occurred while accessing the database for observable apps, waiting for the next retry.")
			time.Sleep(WaitAfterErrorInterval)
		} else {
			for _, app := range *apps {
				if _, ok := observables.Load(app.ID); !ok {
					observables.Store(app.ID, nil)
					go checkGitSourceWrapper(DB, &app, observables)
				}
			}
			break
		}
	}

	for {
		app := <-appChan
		if _, ok := observables.Load(app.ID); !ok {
			observables.Store(app.ID, nil)
			go checkGitSourceWrapper(DB, app, observables)
		}
	}
}

func checkGitSourceWrapper(DB *gorm.DB, app *model.App, observables *sync.Map) {
	contChan := make(chan bool)

	for {
		go checkGitSource(DB, app, contChan, observables)
		cont := <-contChan
		if !cont {
			return
		}
	}
}

func checkGitSource(DB *gorm.DB, app *model.App, contChan chan bool, observables *sync.Map) {
	retryChan := make(chan bool)
	exitChan := make(chan bool)
	for {
		go checkObservable(DB, app, exitChan, retryChan, observables)
		retry := <-retryChan
		if !retry {
			break
		}
	}
	exit := <-exitChan
	if exit {
		contChan <- false
		return
	}

	commitChan := make(chan *string)
	durationChan := make(chan time.Duration)
	var commit *string
	var duration time.Duration
	var restart bool
	for {
		go checkChanges(app.GitSource.RepositoryURL, app.GitSource.Branch, app.GitSource.CommitSHA, commitChan, durationChan)
		commit = <-commitChan
		duration = <-durationChan

		if duration > MaximumDuration {
			duration = MaximumDuration
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
		contChan <- true
		return
	}

	for {
		var errorChan = make(chan bool)

		go deployNewRevision(errorChan, commit)
		hasErr := <-errorChan
		if hasErr {
			zap.L().Error(app.ID + " An error occurred while deploying new revision of %s, waiting for the next retry.")
			time.Sleep(WaitAfterErrorInterval)
		} else {
			break
		}
	}

	zap.L().Info(app.ID + " Deployment of new revision completed, waiting for new changes.")
	time.Sleep(duration)
	contChan <- true
}

func checkObservable(DB *gorm.DB, app *model.App, exitChan chan bool, retryChan chan bool, observables *sync.Map) {
	observableNow, err := datastore.IsObservableApp(DB, app.ID)
	if err != nil {
		zap.L().Error(app.ID + " An error occurred while accessing the database, waiting for the next retry.")
		time.Sleep(WaitAfterErrorInterval)
		retryChan <- true
		return
	}
	if observableNow {
		retryChan <- false
		exitChan <- false
		return
	} else {
		observables.Delete(app.ID)
		retryChan <- false
		exitChan <- true
		return
	}
}

func checkChanges(repoURL string, branch string, currentCommitSHA string, commitChan chan *string, durationChan chan time.Duration) {
	git, err := gitgateway.NewGitGatewayRemote(repoURL)
	if err != nil {
		commitChan <- nil
		durationChan <- -1
		return
	}

	duration, err := git.CommitDuration()
	if err != nil {
		commitChan <- nil
		durationChan <- -1
		return
	}

	checkoutErr := git.Checkout(branch)
	if checkoutErr != nil {
		commitChan <- nil
		durationChan <- -1
		return
	}

	ref, err := git.Head()
	if err != nil {
		commitChan <- nil
		durationChan <- -1
		return
	}

	diff, diffErr := git.Diff(currentCommitSHA, ref)
	if diffErr != nil {
		commitChan <- nil
		durationChan <- -1
		return
	}

	if len(diff) > 0 {
		commitChan <- &ref
	} else {
		commitChan <- nil
	}
	durationChan <- duration
}

// TODO: Integrate with workload controller

/* To integrate with workload controller, replace this function with a functional one.
/  If error occurs, send true to errorChan so that the deployment can be retried.
/  If deployment is successful, return false to errorChan.
/  The commit parameter is reference to HEAD obtained in checkChanges()
/
*/
func deployNewRevision(errorChan chan bool, commit *string) {
	zap.L().Info(*commit + " Deploying new revision...")
	errorChan <- false
}
