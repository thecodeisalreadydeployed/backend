package repositoryobserver

import (
	"fmt"
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

func ObserveGitSources(DB *gorm.DB) {
	aChan := make(chan *model.App)
	var observables sync.Map

	go fetchWrapper(DB, aChan, &observables)

	for {
		app := <-aChan
		go checkWrapper(DB, *app, &observables)
	}
}

func fetchWrapper(DB *gorm.DB, aChan chan *model.App, observables *sync.Map) {
	var wgFetch sync.WaitGroup

	for {
		wgFetch.Add(1)
		go fetchObservableApps(DB, aChan, &wgFetch, observables)
		wgFetch.Wait()
	}
}

func fetchObservableApps(DB *gorm.DB, aChan chan *model.App, wgFetch *sync.WaitGroup, observables *sync.Map) {
	apps, err := datastore.GetObservableApps(DB)

	if err != nil {
		zap.L().Error(err.Error())
		fmt.Println("Unable to fetch observable apps, waiting for the next fetch of observables.")
		time.Sleep(WaitAfterErrorInterval)
		wgFetch.Done()
		return
	}

	if len(*apps) == 0 {
		fmt.Println("All apps are set to not be observed, waiting for the next fetch of observables.")
		time.Sleep(FetchObservableAppsInterval)
		wgFetch.Done()
		return
	}

	for _, app := range *apps {
		_, ok := observables.Load(app.ID)
		if !ok {
			observables.Store(app.ID, nil)
			aChan <- &app
		}
	}

	time.Sleep(FetchObservableAppsInterval)
	wgFetch.Done()
}

func checkWrapper(DB *gorm.DB, app model.App, observables *sync.Map) {
	cChan := make(chan bool)

	for {
		go checkGitSource(DB, app, cChan, observables)
		cont := <-cChan
		if !cont {
			return
		}
	}
}

func checkGitSource(DB *gorm.DB, app model.App, cChan chan bool, observables *sync.Map) {
	commit, duration := checkChanges(app.GitSource.RepositoryURL, app.GitSource.Branch, app.GitSource.CommitSHA)
	if duration > MaximumDuration {
		duration = MaximumDuration
	}
	if commit == nil {
		if duration == -1 {
			fmt.Println("An error occurred while fetching the repository, waiting for next repository check.")
			time.Sleep(WaitAfterErrorInterval)
		} else {
			fmt.Println("There are no changes in the application, waiting for next repository check.")
			time.Sleep(duration)
		}
		cChan <- true
		return
	}

	deployNewRevision()

	time.Sleep(duration)

	retryChan := make(chan bool)
	for {
		go checkObservable(DB, &app, cChan, retryChan, observables)
		cont := <-retryChan
		if !cont {
			return
		}
	}
}

func checkObservable(DB *gorm.DB, app *model.App, cChan chan bool, retryChan chan bool, observables *sync.Map) {
	observableNow, err := datastore.IsObservableApp(DB, app.ID)
	if err != nil {
		time.Sleep(WaitAfterErrorInterval)
		retryChan <- true
		return
	}
	if observableNow {
		retryChan <- false
		cChan <- true
		return
	} else {
		observables.Delete(app.ID)
		retryChan <- false
		cChan <- false
		return
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
	}

	return nil, duration
}

// TODO: Integrate with workload controller

func deployNewRevision() {
	fmt.Println("Deploying new revision...")
	// direct to workload controller
}
