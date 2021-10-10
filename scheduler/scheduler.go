package scheduler

import (
	"github.com/robfig/cron"
)

const (
	_5m  = "0 */5 * * * *"
	_30s = "*/30 * * * * *"
	_15s = "*/15 * * * * *" //nolint
)

func Watch() {
	c := cron.New()

	// Every 5 minutes
	err := c.AddFunc(_5m, func() {
		// repositoryobserver.ObserveGitSources()
	})

	if err != nil {
		panic(err)
	}

	// Every 30 seconds
	err = c.AddFunc(_30s, func() {
		// workloadcontroller.ObserveDeploymentState()
	})

	if err != nil {
		panic(err)
	}

	c.Start()
}
