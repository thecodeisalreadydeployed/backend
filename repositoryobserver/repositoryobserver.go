package repositoryobserver

import (
	"fmt"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
)

const FetchObservableAppsInterval = 3 * time.Minute
const WaitAfterErrorInterval = 10 * time.Second

// TODO: (Critical) Create a goroutine that checks every interval if new apps are made observable

func ObserveGitSources() {
	var observable sync.Map
	var wg sync.WaitGroup

	apps, err := datastore.GetObservableApps(datastore.GetDB())

	if err != nil {
		zap.L().Error(err.Error())
		fmt.Println("Unable to fetch observable apps, waiting for the next fetch of observables.")
		time.Sleep(WaitAfterErrorInterval)
		ObserveGitSources()
		return
	}

	if len(*apps) == 0 {
		fmt.Println("All apps are set to not be observed, waiting for the next fetch of observables.")
		time.Sleep(FetchObservableAppsInterval)
		ObserveGitSources()
		return
	}

	wg.Add(len(*apps))
	for _, app := range *apps {
		go checkGitSource(app, &wg, 0)
	}
	wg.Wait()
	ObserveGitSources()
}

func checkGitSource(app model.App, wg *sync.WaitGroup, waitInterval time.Duration) {
	defer wg.Done()

	commit, duration := checkChanges(app.GitSource.RepositoryURL, app.GitSource.Branch, app.GitSource.CommitSHA)
	if commit == nil {
		if waitInterval == 0 {
			if duration == -1 {
				fmt.Println("An error occurred while fetching the repository, waiting for next repository check.")
				waitInterval = WaitAfterErrorInterval
			}
			fmt.Println("There are no changes in the application, waiting for next repository check.")
			waitInterval = duration
		}
		checkGitSource(app, wg, waitInterval)
		return
	}

	fmt.Printf("appID: %v\n", app.ID)
	fmt.Printf("commit: %v\n", *commit)
	deployNewRevision()

	time.Sleep(waitInterval)

	observable, err := datastore.IsObservableApp(datastore.GetDB(), app.ID)
	if err != nil {
		checkGitSource(app, wg, WaitAfterErrorInterval)
		return
	}
	if observable {
		wg.Add(1)
		go checkGitSource(app, wg, duration)
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
		return nil, duration
	}

	ref, err := git.Head()
	if err != nil {
		return nil, duration
	}

	diff, diffErr := git.Diff(currentCommitSHA, ref)
	if diffErr != nil {
		return nil, duration
	}

	if len(diff) > 0 {
		fmt.Printf("len(diff): %v\n", len(diff))
		return &ref, duration
	}

	return nil, duration
}

// TODO: Integrate with workload controller

func deployNewRevision() {
	// direct to workload controller
}
