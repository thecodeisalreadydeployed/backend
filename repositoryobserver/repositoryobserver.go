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
	for {
		cont := checkGitSource(DB, app, observables)
		if !cont {
			return
		}
	}
}

func checkGitSource(DB *gorm.DB, app *model.App, observables *sync.Map) bool {
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
		return true
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
	return true
}

func checkObservable(DB *gorm.DB, app *model.App, observables *sync.Map) (bool, bool) {
	observableNow, err := datastore.IsObservableApp(DB, app.ID)
	if err != nil {
		zap.L().Error(app.ID + " An error occurred while accessing the database, waiting for the next retry.")
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

	duration, err := git.CommitDuration()
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
