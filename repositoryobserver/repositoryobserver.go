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

const FetchObservableAppsInterval = 3 * time.Minute
const MaximumDuration = 3 * time.Minute
const WaitAfterErrorInterval = 10 * time.Second

var appChan = make(chan *model.App)

func GetAppChannel() chan *model.App {
	return appChan
}

func ObserveGitSources(DB *gorm.DB) {
	//aChan := make(chan *model.App)
	var observables sync.Map

	//go fetchWrapper(DB, aChan, &observables)

	for {
		app := <-appChan
		go checkWrapper(DB, *app, &observables)
	}
}

//
//func fetchWrapper(DB *gorm.DB, aChan chan *model.App, observables *sync.Map) {
//	var wgFetch sync.WaitGroup
//
//	for {
//		wgFetch.Add(1)
//		go fetchObservableApps(DB, aChan, &wgFetch, observables)
//		wgFetch.Wait()
//	}
//}
//
//func fetchObservableApps(DB *gorm.DB, aChan chan *model.App, wgFetch *sync.WaitGroup, observables *sync.Map) {
//	apps, err := datastore.GetObservableApps(DB)
//
//	if err != nil {
//		zap.L().Error(err.Error())
//		zap.L().Info("Unable to fetch observable apps, waiting for the next fetch of observables.")
//		time.Sleep(WaitAfterErrorInterval)
//		wgFetch.Done()
//		return
//	}
//
//	if len(*apps) == 0 {
//		zap.L().Info("All apps are set to not be observed, waiting for the next fetch of observables.")
//		time.Sleep(FetchObservableAppsInterval)
//		wgFetch.Done()
//		return
//	}
//
//	for _, app := range *apps {
//		_, ok := observables.Load(app.ID)
//		if !ok {
//			observables.Store(app.ID, nil)
//			aChan <- &app
//		}
//	}
//
//	time.Sleep(FetchObservableAppsInterval)
//	wgFetch.Done()
//}

func checkWrapper(DB *gorm.DB, app model.App, observables *sync.Map) {
	contChan := make(chan bool)

	for {
		go checkGitSource(DB, app, contChan, observables)
		cont := <-contChan
		if !cont {
			return
		}
	}
}

func checkGitSource(DB *gorm.DB, app model.App, contChan chan bool, observables *sync.Map) {
	retryChan := make(chan bool)
	exitChan := make(chan bool)
	for {
		go checkObservable(DB, &app, exitChan, retryChan, observables)
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
				zap.L().Info(app.ID + " An error occurred while fetching the repository, waiting for the next retry.")
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

	errorChan := make(chan bool)
	for {
		go deployNewRevision(errorChan, commit)
		hasErr := <-errorChan
		if hasErr {
			zap.L().Info(app.ID + " An error occurred while deploying new revision of %s, waiting for the next retry.")
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
		zap.L().Info("An error occurred while accessing the database, waiting for the next retry.")
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

/* Direct to workload controller (can move this function to workload controller module)
/  If error occurs, send true to errorChan so that the deployment can be retried.
/  If deployment is successful, return false to errorChan.
/  The commit parameter is reference to HEAD obtained in checkChanges()
*/
func deployNewRevision(errorChan chan bool, commit string) {
	zap.L().Info("Deploying new revision...")
}
