package repositoryobserver

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"

	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/gitgateway/v2"
	"github.com/thecodeisalreadydeployed/model"
)

const FetchObservableAppsInterval = 3 * time.Minute
const WaitAfterErrorInterval = 10 * time.Second

func ObserveGitSources(DB *gorm.DB) {
	aChan := make(chan *model.App)

	go checkObservableApps(DB, aChan, true)

	for {
		select {
		case app := <-aChan:
			go checkGitSource(*app, 0, true)
		}
	}
}

func checkObservableApps(DB *gorm.DB, aChan chan *model.App, recurrent bool) {
	apps, err := datastore.GetObservableApps(DB)

	if err != nil {
		zap.L().Error(err.Error())
		fmt.Println("Unable to fetch observable apps, waiting for the next fetch of observables.")
		time.Sleep(WaitAfterErrorInterval)
		defer func() {
			go checkObservableApps(DB, aChan, recurrent)
		}()
		return
	}

	if len(*apps) == 0 {
		fmt.Println("All apps are set to not be observed, waiting for the next fetch of observables.")
		time.Sleep(FetchObservableAppsInterval)
		defer func() {
			go checkObservableApps(DB, aChan, recurrent)
		}()
		return
	}

	for _, app := range *apps {
		aChan <- &app
	}

	if recurrent {
		time.Sleep(FetchObservableAppsInterval)
		defer func() {
			go checkObservableApps(DB, aChan, true)
		}()
	}
}

func checkGitSource(app model.App, waitInterval time.Duration, recurrent bool) {
	commit, duration := checkChanges(app.GitSource.RepositoryURL, app.GitSource.Branch, app.GitSource.CommitSHA)
	if commit == nil {
		if waitInterval == 0 {
			if duration == -1 {
				fmt.Println("An error occurred while fetching the repository, waiting for next repository check.")
				waitInterval = WaitAfterErrorInterval
			} else {
				fmt.Println("There are no changes in the application, waiting for next repository check.")
				waitInterval = duration
			}
		}
		time.Sleep(waitInterval)
		defer func() {
			go checkGitSource(app, waitInterval, recurrent)
		}()
		return
	}

	deployNewRevision()

	if recurrent {
		time.Sleep(waitInterval)

		observable, err := datastore.IsObservableApp(datastore.GetDB(), app.ID)
		if err != nil {
			defer func() {
				go checkGitSource(app, WaitAfterErrorInterval, true)
			}()
			return
		}
		if observable {
			defer func() {
				go checkGitSource(app, duration, true)
			}()
		}
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
		return &ref, duration
	}

	return nil, duration
}

// TODO: Integrate with workload controller

func deployNewRevision() {
	fmt.Println("Deploying new revision...")
	// direct to workload controller
}
