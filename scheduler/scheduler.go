package scheduler

import (
	"github.com/robfig/cron"
	"github.com/thecodeisalreadydeployed/repositoryobserver"
	"github.com/thecodeisalreadydeployed/workloadcontroller"
)

func Watch() {
	c := cron.New()

	// Every 5 minutes
	err := c.AddFunc("0 */5 * * * *", func() {
		repositoryobserver.ObserveGitSource()
	})

	if err != nil {
		panic(err)
	}

	// Every 30 seconds
	err = c.AddFunc("*/30 * * * * *", func() {
		workloadcontroller.ObserveDeploymentState()
	})

	if err != nil {
		panic(err)
	}

	c.Start()
}
