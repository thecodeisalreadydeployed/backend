package scheduler

import (
	"github.com/robfig/cron"
	"github.com/thecodeisalreadydeployed/repositoryobserver"
	"github.com/thecodeisalreadydeployed/workloadcontroller"
)

const (
	_5m  = "0 */5 * * * *"
	_30s = "*/30 * * * * *"
	_15s = "*/15 * * * * *"
)

func Watch() {
	c := cron.New()

	// Every 5 minutes
	err := c.AddFunc(_5m, func() {
		repositoryobserver.ObserveGitSource()
	})

	if err != nil {
		panic(err)
	}

	// Every 30 seconds
	err = c.AddFunc(_30s, func() {
		workloadcontroller.ObserveDeploymentState()
	})

	if err != nil {
		panic(err)
	}

	c.Start()
}
